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
	"github.com/gorilla/websocket"
)

func StartWebServer(hashtags []string, tweets <-chan *twitter.Tweet) {

	env := webEnv{
		boardWebsocketClients:      make(map[*websocket.Conn]bool),
		controllerWebsocketClients: make(map[*websocket.Conn]bool),
		tweetsForConsideration:     make(map[string]TweetSummary),
		blockedUsers:               make(map[string]bool),
		tweets:                     tweets,
	}

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		var blockedUsers []string
		for user := range env.blockedUsers {
			blockedUsers = append(blockedUsers, user)
		}

		jsonHashtags, err := json.Marshal(struct {
			Hashtags     []string `json:"hashtags"`
			BlockedUsers []string `json:"blockedUsers"`
		}{
			Hashtags:     hashtags,
			BlockedUsers: blockedUsers,
		})

		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err.Error())
		}

		fmt.Fprint(w, string(jsonHashtags))
	})

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/control-ws", env.controllerWebsocketHandler)

	http.HandleFunc("/board-ws", env.boardWebsocketHandler)

	go env.handleTweetsFromTwitter(tweets)

	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}
