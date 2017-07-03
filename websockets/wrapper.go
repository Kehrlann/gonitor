package websockets

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

type Connection interface {
	WriteMessage(message string) error
}

// NewWebsocketConnection creates a wrapper for a websocket connection, with a simple write method.
func NewWebsocketConnection(conn *websocket.Conn) Connection {
	wrapper := &websocketConnectionWrapper{conn, make(chan string), false, sync.Mutex{}}
	wrapper.startReadPump() // This detects "close" control messages, which otherwise go unnoticed
	wrapper.startWritePump()
	return wrapper
}

type websocketConnectionWrapper struct {
	c        *websocket.Conn
	message  chan string
	isClosed bool
	lock     sync.Mutex	// Seriously, a mutex for a boolean write ? -_-
}

// WriteMessage schedules an asynchronous write to the underlying connection. It raises the error if said connection
// is already marked as closed.
// However, if there are no errors, there is not guarantee that the message will be transmitted, as the close detection
// mechanism is also asynchronous.
func (connection *websocketConnectionWrapper) WriteMessage(message string) error {
	connection.lock.Lock()
	if connection.isClosed {
		connection.lock.Unlock()
		return errors.New("trying to write on a closed connection")
	}
	connection.lock.Unlock()
	connection.message <- message
	return nil
}

func (connection *websocketConnectionWrapper) startWritePump() {
	go func() {
		for message := range connection.message {
			err := connection.c.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				connection.lock.Lock()
				connection.isClosed = true
				connection.lock.Unlock()
				return
			}
		}
	}()
}

func (connection *websocketConnectionWrapper) startReadPump() {
	go func() {
		for {
			_, _, err := connection.c.ReadMessage()
			if err != nil {
				return
			}
		}
	}()
}
