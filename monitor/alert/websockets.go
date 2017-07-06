package alert

import (
	"sync"
	"encoding/json"

	"github.com/kehrlann/gonitor/monitor"
	"github.com/kehrlann/gonitor/websockets"
)

type WebsocketsEmitter struct {
	websocketConnections map[uint]websockets.Connection
	currentIndex         uint
	lock                 sync.Mutex
}

func (emitter *WebsocketsEmitter) Emit(message *monitor.StateChangeMessage) {
	jsonMessage, marshalErr := json.Marshal(message)

	if marshalErr != nil {
		// TODO : try to trigger an error when marshalling
		//		what should we do ?
		// 		Probaly log and return
	}

	connections := emitter.getConnections()
	for key, conn := range connections {
		writeErr := conn.WriteMessage(string(jsonMessage))
		if writeErr != nil {
			// Usually when a connection is closed, you get an error (forcefully closed or with a CloseMessage)
			emitter.unregisterConnection(key)
		}
	}
}

// NewWebsocketEmitter creates an emitter that listens for incoming websocket connections,
// and broadcasts alerts on all incoming connections.
// All closed connections will be automatically remove
func NewWebsocketEmitter(receive_conn <-chan websockets.Connection) *WebsocketsEmitter {
	res := &WebsocketsEmitter{websocketConnections: make(map[uint]websockets.Connection)}
	go func() {
		for conn := range receive_conn {
			res.registerConnection(conn)
		}
	}()
	return res
}

func (emitter *WebsocketsEmitter) registerConnection(connection websockets.Connection) {
	emitter.lock.Lock()
	defer emitter.lock.Unlock()
	connectionIndex := emitter.currentIndex
	emitter.websocketConnections[connectionIndex] = connection
	emitter.currentIndex++
}

func (emitter *WebsocketsEmitter) unregisterConnection(key uint) {
	emitter.lock.Lock()
	defer emitter.lock.Unlock()
	delete(emitter.websocketConnections, key)
}

func (emitter *WebsocketsEmitter) getConnections() map[uint]websockets.Connection {
	emitter.lock.Lock()
	defer emitter.lock.Unlock()
	new_map := make(map[uint]websockets.Connection)

	for key, value := range emitter.websocketConnections {
		new_map[key] = value
	}

	return new_map
}
