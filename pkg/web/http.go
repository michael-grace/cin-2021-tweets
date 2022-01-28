/**
URY Tweet Board

Author: Michael Grace <michael.grace@ury.org.uk>
*/

package web

import (
	"fmt"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
)

func StartWebServer(tweets <-chan *twitter.Tweet) {
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "#CIN22")
	})

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/control-ws", func(w http.ResponseWriter, r *http.Request) {
		ControllerWebsocketMaster.websocketHandler(w, r, tweets)
	})

	http.HandleFunc("/board-ws", func(w http.ResponseWriter, r *http.Request) {
		BoardWebsocketMaster.websocketHandler(w, r)
	})

	go ControllerWebsocketMaster.HandleTweetsFromTwitter(tweets)

	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}
