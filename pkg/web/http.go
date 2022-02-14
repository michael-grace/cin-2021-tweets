/**
URY Tweet Board

Author: Michael Grace <michael.grace@ury.org.uk>
*/

package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/michael-grace/cin-2021-tweets/pkg/logging"
)

func StartWebServer(hashtags []string, tweets <-chan *twitter.Tweet) {

	env := webEnv{
		boardWebsocketClients:      make(map[*websocket.Conn]bool),
		controllerWebsocketClients: make(map[*websocket.Conn]bool),
		tweetsForConsideration:     make(map[string]TweetSummary),
		blockedUsers:               make(map[string]bool),
		tweets:                     tweets,
		recentlySentToBoard:        make(chan *TweetSummary, 8),
		boardTweetsForQuerying:     make(map[TweetSummary]bool),
		wsAuthToken:                uuid.New(),
	}

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		jsonData, err := json.Marshal(struct {
			Hashtags []string `json:"hashtags"`
		}{
			Hashtags: hashtags,
		})

		if err != nil {
			w.WriteHeader(500)
			logging.Error(err)
			fmt.Fprint(w, err.Error())
		}

		fmt.Fprint(w, string(jsonData))
	})

	http.Handle("/ws-auth", authHandler{wsAuthHandler{
		authToken: env.wsAuthToken,
	}})

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.Handle("/control/", authHandler{handler: fs})

	http.HandleFunc("/control-ws", env.controllerWebsocketHandler)

	http.HandleFunc("/board-ws", env.boardWebsocketHandler)

	go env.handleTweetsFromTwitter(tweets)

	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}
