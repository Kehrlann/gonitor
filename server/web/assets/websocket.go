package assets

import (
	"net/http"
	"github.com/gorilla/websocket"
	"time"
)

func HandleWebsockets(response http.ResponseWriter, request *http.Request) {
	up := websocket.Upgrader{WriteBufferSize:1024}
	conn, err := up.Upgrade(response, request, nil)

	if err != nil {
		// TODO : how can the Upgrade fail ?
	}

	conn.CloseHandler()

	for range time.Tick(1 * time.Second){
		conn.WriteMessage(websocket.TextMessage, []byte("{ \"date\" : \"" + time.Now().String() + "\"}"))
	}
}
