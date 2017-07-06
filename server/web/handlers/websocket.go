package handlers

import (
	"net/http"
	"fmt"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/kehrlann/gonitor/websockets"
)

type WebsocketHandler struct {
	IncomingConnections chan<- websockets.Connection
}

func (handler WebsocketHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	up := websocket.Upgrader{WriteBufferSize:1024}
	conn, err := up.Upgrade(response, request, nil)	// This closes the response writer
	if err != nil {
		log.Error(fmt.Sprintf("Error establishing websocket connection : %v", err))
		return
	}
	wrapper := websockets.NewWebsocketConnection(conn)
	handler.IncomingConnections <- wrapper
}
