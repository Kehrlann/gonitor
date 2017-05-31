package alert

import (
	"github.com/gorilla/websocket"
	"sync"
	"encoding/json"
)

type WebsocketsEmitter struct {
	websocketConnections map[int]websocket.Conn
	currentIndex int
	lock sync.Mutex
}

func (emitter *WebsocketsEmitter) Emit(message *StateChangeMessage) {
	for _, conn := range emitter.websocketConnections {
		// TODO : try to trigger an error when marshalling
		jsonMessage, _ := json.Marshal(message)

		newErr := conn.WriteMessage(websocket.TextMessage, jsonMessage)

		if newErr != nil {
			// TODO : what happends when error ?
		}
	}
}


func (emitter *WebsocketsEmitter) RegisterConnection(connection websocket.Conn) {
	emitter.lock.Lock()
	defer emitter.lock.Unlock()
	emitter.websocketConnections[emitter.currentIndex] = connection
	emitter.currentIndex++
	// TODO : how do we unregister this ?? -> see connection.SetCloseHandler()
}