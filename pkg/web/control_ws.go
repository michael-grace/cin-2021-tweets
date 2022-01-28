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

type controlWebsocketH struct {
	client                 *websocket.Conn
	tweetsForConsideration map[string]TweetSummary
	blockedUsers           []string
}

var ControllerWebsocketMaster controlWebsocketH = controlWebsocketH{tweetsForConsideration: make(map[string]TweetSummary)}

func (h *controlWebsocketH) websocketHandler(w http.ResponseWriter, r *http.Request, tweets <-chan *twitter.Tweet) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to generate upgrader: %s", err)
		return
	}

	h.client = ws

	defer func() {
		h.client = nil
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

		var decision ControllerDecision
		json.Unmarshal(message, &decision)

		fmt.Printf("%v - %v", decision.Decision, h.tweetsForConsideration[decision.ID])

		switch decision.Decision {
		case "BLOCK":
			h.blockedUsers = append(h.blockedUsers, h.tweetsForConsideration[decision.ID].User)

		case "ACCEPT":
			tweeet := h.tweetsForConsideration[decision.ID]

			embed, err := http.Get(fmt.Sprintf("https://publish.twitter.com/oembed?url=https://twitter.com/%s/status/%s&hide_thread=true&theme=dark&hide_media=true", tweeet.User, tweeet.ID))

			if err != nil {
				fmt.Println(err)
			}

			defer embed.Body.Close()

			j, err := io.ReadAll(embed.Body)

			if err != nil {
				fmt.Println(err)
			}

			var embedJson EmbedJson
			json.Unmarshal(j, &embedJson)

			enc := base64.StdEncoding.EncodeToString([]byte(embedJson.HTML))

			tweeet.TweetHTML = enc

			BoardWebsocketMaster.SendTweet(tweeet)

		}

		delete(h.tweetsForConsideration, decision.ID)

	}
}

func (h *controlWebsocketH) HandleTweetsFromTwitter(tweets <-chan *twitter.Tweet) {
	for tweet := range tweets {
		if h.client != nil {
			var blocked bool

			for _, blockedUser := range h.blockedUsers {
				if tweet.User.ScreenName == blockedUser {
					blocked = true
					break
				}
			}

			if blocked {
				continue
			}

			tweetSummary := TweetSummary{
				ID:      tweet.IDStr,
				Name:    tweet.User.Name,
				User:    tweet.User.ScreenName,
				Message: tweet.Text,
			}

			h.tweetsForConsideration[tweetSummary.ID] = tweetSummary

			err := h.client.WriteJSON(tweetSummary)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}
