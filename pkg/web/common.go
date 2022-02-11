/**
URY Tweet Board

Author: Michael Grace <michael.grace@ury.org.uk>
*/

package web

import (
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/gorilla/websocket"
)

type webEnv struct {
	boardWebsocketClients      map[*websocket.Conn]bool
	controllerWebsocketClients map[*websocket.Conn]bool
	tweetsForConsideration     map[string]TweetSummary
	blockedUsers               map[string]bool
	tweets                     <-chan *twitter.Tweet
	recentlySentToBoard        chan *TweetSummary
}
type TweetSummary struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	User      string `json:"user"`
	Message   string `json:"tweet"`
	TweetHTML string `json:"html"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO
	},
}
