package alert

import (
	"github.com/gorilla/websocket"
	"sync"
	"encoding/json"
	"github.com/kehrlann/gonitor/monitor"
)

type WebsocketsEmitter struct {
	websocketConnections map[uint]*websocket.Conn
	currentIndex         uint
	lock                 sync.RWMutex
}

func (emitter *WebsocketsEmitter) Emit(message *monitor.StateChangeMessage) {
	jsonMessage, marshalErr := json.Marshal(message)

	if marshalErr != nil {
		// TODO : try to trigger an error when marshalling
		//		what should we do ?
	}

	for key, conn := range emitter.websocketConnections {

		writeErr := conn.WriteMessage(websocket.TextMessage, jsonMessage)

		if writeErr != nil {
			// TODO : Unregister when error ... but what about error cause ?
			// 			should we check for timeout ? Closed connection ?
			emitter.unregisterConnection(key)
		}
	}
}

// NewWebsocketEmitter creates an emitter that listens for incoming websocket connections,
// and broadcasts alerts on all incoming connections.
// All closed connections will be automatically remove
func NewWebsocketEmitter(receive_conn <-chan *websocket.Conn) *WebsocketsEmitter {
	res := &WebsocketsEmitter{websocketConnections: make(map[uint]*websocket.Conn)}
	go func() {
		for conn := range receive_conn {
			res.registerConnection(conn)
		}
	}()
	return res
}

func (emitter *WebsocketsEmitter) registerConnection(connection *websocket.Conn) {
	emitter.lock.Lock()
	defer emitter.lock.Unlock()
	emitter.websocketConnections[emitter.currentIndex] = connection
	emitter.currentIndex++

	// TODO : can't do this because it breaks tests (the underlying connection being null)
	//go func() {
	//	for {
			// See closed messages !
			// connection.SetReadDeadline(time.Now().Add(10 * time.Second))
			// connection.ReadMessage()
	//	}
	//}()
}

func (emitter *WebsocketsEmitter) unregisterConnection(key uint) {
	emitter.lock.Lock()
	defer emitter.lock.Unlock()
	delete(emitter.websocketConnections, key)
}

func (emitter *WebsocketsEmitter) getConnections() map[uint]*websocket.Conn {
	emitter.lock.RLock()
	defer emitter.lock.RUnlock()
	new_map := make(map[uint]*websocket.Conn)

	for key, value := range emitter.websocketConnections {
		new_map[key] = value
	}

	return new_map
}
