/**
URY Tweet Board

Author: Michael Grace <michael.grace@ury.org.uk>
*/

package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type boardWebsocketH struct {
	clients map[*websocket.Conn]bool
}

var BoardWebsocketMaster boardWebsocketH = boardWebsocketH{clients: make(map[*websocket.Conn]bool)}

func (h *boardWebsocketH) websocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to generate upgrader: %s", err)
		return
	}

	h.clients[ws] = true

	defer func() {
		delete(h.clients, ws)
		ws.Close()
	}()

	ws.ReadMessage()

}

func (h *boardWebsocketH) SendTweet(tweet TweetSummary) {
	for client := range h.clients {
		client.WriteJSON(tweet)
	}
}
