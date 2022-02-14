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
)

type BlockStatus string

const (
	BlockUser   BlockStatus = "BLOCK"
	UnblockUser BlockStatus = "UNBLOCK"
)

func (h *webEnv) sendJSONToControllers(data interface{}) {
	for client := range h.controllerWebsocketClients {
		if err := client.WriteJSON(data); err != nil {
			fmt.Println(err)
		}
	}
}

func (h *webEnv) changeBlockStatus(blockStatus BlockStatus, user string) {
	h.sendJSONToControllers(struct {
		Action string `json:"action"`
		User   string `json:"user"`
	}{
		Action: string(blockStatus),
		User:   user,
	})
}

func (h *webEnv) removeTweetFromConsideration(id string) {
	delete(h.tweetsForConsideration, id)
	h.sendJSONToControllers(struct {
		Action string `json:"action"`
		ID     string `json:"id"`
	}{
		Action: "REMOVE",
		ID:     id,
	})
}

func (h *webEnv) controllerWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to generate upgrader: %s", err)
		return
	}

	defer ws.Close()

	var wsAuthenticated bool

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

		var messageContent struct {
			Action  string `json:"action"`
			Content string `json:"content"`
		}

		if err = json.Unmarshal(message, &messageContent); err != nil {
			fmt.Println(err.Error())
		}

		if !wsAuthenticated && messageContent.Action != "AUTH" {
			ws.WriteJSON(struct {
				Action string `json:"action"`
				OK     bool   `json:"ok"`
			}{
				Action: "AUTH",
				OK:     false,
			})
			continue
		}

		switch messageContent.Action {
		case "AUTH":
			if messageContent.Content == h.wsAuthToken.String() {
				ws.WriteJSON(struct {
					Action string `json:"action"`
					OK     bool   `json:"ok"`
				}{
					Action: "AUTH",
					OK:     true,
				})
				wsAuthenticated = true

				h.controllerWebsocketClients[ws] = true

				defer func() {
					delete(h.controllerWebsocketClients, ws)
				}()
			} else {
				ws.WriteJSON(struct {
					Action string `json:"action"`
					OK     bool   `json:"ok"`
				}{
					Action: "AUTH",
					OK:     false,
				})
			}
		case "CLEAR_CONTROL":
			h.tweetsForConsideration = make(map[string]TweetSummary)
			h.sendJSONToControllers(struct {
				Action string `json:"action"`
			}{
				Action: "CLEAR_CONTROL",
			})

		case "CLEAR_BOARD":
			h.sendJSONToBoard(struct {
				Action string `json:"action"`
			}{
				Action: "CLEAR",
			})
			h.recentlySentToBoard = make(chan *TweetSummary, 8)
			h.sendJSONToControllers(struct {
				Action string `json:"action"`
			}{
				Action: "CLEAR_BOARD",
			})

		case "BOARD_REMOVE":
			h.sendJSONToControllers(struct {
				Action string `json:"action"`
				ID     string `json:"id"`
			}{
				Action: "UNRECENT",
				ID:     messageContent.Content,
			})

			h.sendJSONToBoard(struct {
				Action string `json:"action"`
				ID     string `json:"id"`
			}{
				Action: "REMOVE",
				ID:     messageContent.Content,
			})

		case "UNBLOCK":
			delete(h.blockedUsers, messageContent.Content)
			h.changeBlockStatus(UnblockUser, messageContent.Content)

		case "BLOCK":
			user := h.tweetsForConsideration[messageContent.Content].User
			h.blockedUsers[user] = true
			h.changeBlockStatus(BlockUser, user)
			h.removeTweetFromConsideration(messageContent.Content)

		case "REJECT":
			h.removeTweetFromConsideration(messageContent.Content)

		case "ACCEPT":
			tweet := h.tweetsForConsideration[messageContent.Content]
			embed, err := http.Get(
				fmt.Sprintf(
					"https://publish.twitter.com/oembed?url=https://twitter.com/%s/status/%s&hide_thread=true&theme=light",
					tweet.User,
					tweet.ID))

			if err != nil {
				fmt.Println(err.Error())
			}

			j, err := io.ReadAll(embed.Body)
			if err != nil {
				fmt.Println(err.Error())
			}

			var embedJson struct {
				HTML string `json:"html"`
			}

			json.Unmarshal(j, &embedJson)

			tweet.TweetHTML = base64.StdEncoding.EncodeToString([]byte(embedJson.HTML))

			h.sendTweet(tweet)

			if len(h.recentlySentToBoard) == cap(h.recentlySentToBoard) {
				// Tell controllers tweet no longer recent
				oldTweet := <-h.recentlySentToBoard
				delete(h.boardTweetsForQuerying, *oldTweet)
				h.sendJSONToControllers(struct {
					Action string `json:"action"`
					ID     string `json:"id"`
				}{
					Action: "UNRECENT",
					ID:     oldTweet.ID,
				})
			}

			// Tell controllers about new tweet
			h.recentlySentToBoard <- &tweet
			h.boardTweetsForQuerying[tweet] = true
			h.sendJSONToControllers(struct {
				Action string       `json:"action"`
				Tweet  TweetSummary `json:"tweet"`
			}{
				Action: "RECENT",
				Tweet:  tweet,
			})

			h.removeTweetFromConsideration(messageContent.Content)

		case "QUERY":
			// Blocked Users
			for user := range h.blockedUsers {
				if err := ws.WriteJSON(struct {
					Action string `json:"action"`
					User   string `json:"user"`
				}{
					Action: "BLOCK",
					User:   user,
				}); err != nil {
					fmt.Println(err.Error())
				}
			}

			// Recent Tweets
			for tweet := range h.boardTweetsForQuerying {
				if err := ws.WriteJSON(struct {
					Action string       `json:"action"`
					Tweet  TweetSummary `json:"tweet"`
				}{
					Action: "RECENT",
					Tweet:  tweet,
				}); err != nil {
					fmt.Println(err.Error())
				}
			}

			// Pending Tweets
			for _, tweet := range h.tweetsForConsideration {
				if err := ws.WriteJSON(struct {
					Action string       `json:"action"`
					Tweet  TweetSummary `json:"tweet"`
				}{
					Action: "CONSIDER",
					Tweet:  tweet,
				}); err != nil {
					fmt.Println(err.Error())
				}
			}
		}

	}
}

func (h *webEnv) handleTweetsFromTwitter(tweets <-chan *twitter.Tweet) {
	for tweet := range tweets {
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

		h.sendJSONToControllers(struct {
			Action string       `json:"action"`
			Tweet  TweetSummary `json:"tweet"`
		}{
			Action: "CONSIDER",
			Tweet:  tweetSummary,
		})

	}
}
