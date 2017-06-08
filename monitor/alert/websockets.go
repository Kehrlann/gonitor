package alert

import (
	"github.com/gorilla/websocket"
	"sync"
	"encoding/json"
	"github.com/kehrlann/gonitor/monitor"
)

type WebsocketsEmitter struct {
	websocketConnections map[uint]*websocket.Conn
	currentIndex uint
	lock sync.Mutex
}

func (emitter *WebsocketsEmitter) Emit(message *monitor.StateChangeMessage) {
	for key, conn := range emitter.websocketConnections {
		jsonMessage, marshalErr := json.Marshal(message)

		if marshalErr != nil {
			// TODO : try to trigger an error when marshalling
			//		what should we do ?
		}

		writeErr := conn.WriteMessage(websocket.TextMessage, jsonMessage)

		if writeErr != nil {
			// TODO : Unregister when error ... but what about error cause ?
			// 			should we check for timeout ? Closed connection ?
			emitter.UnregisterConnection(key)
		}
	}
}

func NewWebsocketEmitter(receive_conn <-chan *websocket.Conn) *WebsocketsEmitter {
	res := &WebsocketsEmitter{websocketConnections:make(map[uint]*websocket.Conn)}
	go func() {
		for conn := range receive_conn {
			res.RegisterConnection(conn)
		}
	}()
	return res
}

func (emitter *WebsocketsEmitter) RegisterConnection(connection *websocket.Conn) {
	emitter.lock.Lock()
	defer emitter.lock.Unlock()
	emitter.websocketConnections[emitter.currentIndex] = connection
	emitter.currentIndex++
}

func (emitter *WebsocketsEmitter) UnregisterConnection(key uint) {
	emitter.lock.Lock()
	defer emitter.lock.Unlock()
	delete(emitter.websocketConnections, key)
}