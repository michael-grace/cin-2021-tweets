/**
URY Tweet Board

Author: Michael Grace <michael.grace@ury.org.uk>
*/

package web

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/gorilla/websocket"
)

type BlockStatus string

const (
	BlockUser   BlockStatus = "BLOCK"
	UnblockUser BlockStatus = "UNBLOCK"
)

func (h *webEnv) changeBlockStatus(blockStatus BlockStatus, user string) {
	for client := range h.controllerWebsocketClients {
		if err := client.WriteJSON(struct {
			Action string `json:"action"`
			User   string `json:"user"`
		}{
			Action: string(blockStatus),
			User:   user,
		}); err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (h *webEnv) controllerWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to generate upgrader: %s", err)
		return
	}

	h.controllerWebsocketClients[ws] = true

	defer func() {
		delete(h.controllerWebsocketClients, ws)
		ws.Close()
	}()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if !strings.Contains(err.Error(), "1001") {
				log.Printf("Failed to read WebSocket message: %s", err)
			} else {
				// Client Disconnected
				return
			}
		}

		switch string(message) {
		case "CLEAR_CONTROL":
			h.tweetsForConsideration = make(map[string]TweetSummary)
			for client := range h.controllerWebsocketClients {
				err = client.WriteMessage(websocket.TextMessage, []byte("CLEAR"))

				if err != nil {
					fmt.Println(err)
				}
			}

		case "CLEAR_BOARD":
			h.sendTextMessageToBoard("CLEAR")

		default:

			var decision struct {
				ID       string `json:"id"`
				Decision string `json:"decision"`
			}

			json.Unmarshal(message, &decision)

			if decision.Decision == "UNBLOCK" {
				delete(h.blockedUsers, decision.ID)
				h.changeBlockStatus(UnblockUser, decision.ID)
				break
			}

			// If we get to here, we're dealing with a tweet that
			// is under consideration
			fmt.Printf("%v - %v", decision.Decision, h.tweetsForConsideration[decision.ID])

			for client := range h.controllerWebsocketClients {
				err = client.WriteJSON(struct {
					Action string `json:"action"`
					ID     string `json:"id"`
				}{
					Action: "REMOVE",
					ID:     decision.ID,
				})

				if err != nil {
					fmt.Println(err)
				}
			}

			switch decision.Decision {
			case "BLOCK":
				h.blockedUsers[h.tweetsForConsideration[decision.ID].User] = true
				h.changeBlockStatus(BlockUser, h.tweetsForConsideration[decision.ID].User)

			case "ACCEPT":
				tweeet := h.tweetsForConsideration[decision.ID]

				embed, err := http.Get(
					fmt.Sprintf(
						"https://publish.twitter.com/oembed?url=https://twitter.com/%s/status/%s&hide_thread=true&theme=light&hide_media=true",
						tweeet.User,
						tweeet.ID))

				if err != nil {
					fmt.Println(err)
				}

				defer embed.Body.Close()

				j, err := io.ReadAll(embed.Body)

				if err != nil {
					fmt.Println(err)
				}

				var embedJson struct {
					HTML string `json:"html"`
				}

				json.Unmarshal(j, &embedJson)

				enc := base64.StdEncoding.EncodeToString([]byte(embedJson.HTML))

				tweeet.TweetHTML = enc

				h.sendTweet(tweeet)

			}

			delete(h.tweetsForConsideration, decision.ID)
		}
	}
}

func (h *webEnv) handleTweetsFromTwitter(tweets <-chan *twitter.Tweet) {
	for tweet := range tweets {
		fmt.Println(tweet)
		if _, blocked := h.blockedUsers[tweet.User.ScreenName]; blocked {
			continue
		}

		tweetSummary := TweetSummary{
			ID:      tweet.IDStr,
			Name:    tweet.User.Name,
			User:    tweet.User.ScreenName,
			Message: tweet.Text,
		}

		h.tweetsForConsideration[tweetSummary.ID] = tweetSummary

		for client := range h.controllerWebsocketClients {
			err := client.WriteJSON(struct {
				TweetSummary
				Action string `json:"action"`
			}{
				TweetSummary: tweetSummary,
				Action:       "CONSIDER",
			})

			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}
