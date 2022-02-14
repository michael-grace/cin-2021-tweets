/**
URY Tweet Board

Author: Michael Grace <michael.grace@ury.org.uk>
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
