/**
URY Tweet Board

Author: Michael Grace <michael.grace@ury.org.uk>
*/

package web

import (
	"net/http"

	"github.com/gorilla/websocket"
)

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
