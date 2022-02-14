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
