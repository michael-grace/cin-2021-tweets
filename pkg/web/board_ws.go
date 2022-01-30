/**
URY Tweet Board

Author: Michael Grace <michael.grace@ury.org.uk>
*/

package web

import (
	"fmt"
	"net/http"
)

func (h *webEnv) boardWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to generate upgrader: %s", err)
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
		client.WriteJSON(tweet)
	}
}
