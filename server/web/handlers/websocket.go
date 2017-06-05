package handlers

import (
	"net/http"
	"time"
	"fmt"

	"github.com/gorilla/websocket"
	log "github.com/Sirupsen/logrus"
)

func HandleWebsockets(response http.ResponseWriter, request *http.Request) {
	up := websocket.Upgrader{WriteBufferSize:1024}
	conn, err := up.Upgrade(response, request, nil)	// This closes the response writer
	if err != nil {
		log.Error(fmt.Sprintf("Error establishing websocket connection : %v", err))
		return
	}

	for range time.Tick(1 * time.Second){
		// TODO : you need to read the message to be able to know whether the connection was closed or not
		// TODO : set deadline is the worst fucking API ever :)
		conn.SetReadDeadline(time.Now().Add(time.Second))
		conn.ReadMessage()
		//fmt.Println(messageType, message, error)

		writeError := conn.WriteMessage(websocket.TextMessage, []byte("{ \"date\" : \"" + time.Now().String() + "\"}"))
		fmt.Println("WRITE")
		if writeError != nil {
			fmt.Println(writeError)
			// TODO : close the connection and unregister
		} else {
			fmt.Println("SUCCESS")
		}

	}
}
