/**
URY Tweet Board
Copyright (C) 2022 Michael Grace <michael.grace@ury.org.uk>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
