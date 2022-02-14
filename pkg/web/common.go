/**
URY Tweet Board

Author: Michael Grace <michael.grace@ury.org.uk>
*/

package web

import (
	"fmt"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type webEnv struct {
	boardWebsocketClients      map[*websocket.Conn]bool
	controllerWebsocketClients map[*websocket.Conn]bool
	tweetsForConsideration     map[string]TweetSummary
	blockedUsers               map[string]bool
	tweets                     <-chan *twitter.Tweet
	recentlySentToBoard        chan *TweetSummary
	boardTweetsForQuerying     map[TweetSummary]bool
	wsAuthToken                uuid.UUID
}
type TweetSummary struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	User      string `json:"user"`
	Message   string `json:"tweet"`
	TweetHTML string `json:"html"`
}

func (t TweetSummary) String() string {
	return fmt.Sprintf("@%s: %s", t.User, t.Message)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO
	},
}
