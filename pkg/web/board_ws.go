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
	"net/http"

	"github.com/michael-grace/cin-2021-tweets/pkg/logging"
)

func (h *webEnv) boardWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.Error(err)
		return
	}

	h.boardWebsocketClients[ws] = true

	defer func() {
		delete(h.boardWebsocketClients, ws)
		ws.Close()
	}()

	ws.ReadMessage()

}

func (h *webEnv) sendTweet(tweet TweetSummary) {
	for client := range h.boardWebsocketClients {
		if err := client.WriteJSON(struct {
			Action string       `json:"action"`
			Tweet  TweetSummary `json:"tweet"`
		}{
			Action: "ADD",
			Tweet:  tweet,
		}); err != nil {
			logging.Error(err)
		}
	}
}

func (h *webEnv) sendJSONToBoard(data interface{}) {
	for client := range h.boardWebsocketClients {
		if err := client.WriteJSON(data); err != nil {
			logging.Error(err)
		}
	}
}
