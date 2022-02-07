/**
URY Tweet Board

Author: Michael Grace <michael.grace@ury.org.uk>
*/

package web

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/gorilla/websocket"
)

type webEnv struct {
	boardWebsocketClients      map[*websocket.Conn]bool
	controllerWebsocketClients map[*websocket.Conn]bool
	tweetsForConsideration     map[string]TweetSummary
	blockedUsers               map[string]bool
}

func StartWebServer(tweets <-chan *twitter.Tweet) {

	env := webEnv{
		boardWebsocketClients:      make(map[*websocket.Conn]bool),
		controllerWebsocketClients: make(map[*websocket.Conn]bool),
		tweetsForConsideration:     make(map[string]TweetSummary),
		blockedUsers:               make(map[string]bool),
	}

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, os.Getenv("HASHTAG"))
	})

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/control-ws", func(w http.ResponseWriter, r *http.Request) {
		env.controllerWebsocketHandler(w, r, tweets)
	})

	http.HandleFunc("/board-ws", env.boardWebsocketHandler)

	go env.handleTweetsFromTwitter(tweets)

	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}
