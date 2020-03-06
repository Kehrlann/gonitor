package websockets

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http/httptest"
	"strings"
	"net/http"
	"github.com/gorilla/websocket"
	"time"
	"errors"
	"sync"
)

var _ = Describe("Websockets integrations tests -> ", func() {

	var server *testServer
	//message := monitor.RecoveryMessage(config.Resource{}, []int{})

	BeforeEach(func() {
		server = newServer()
	})

	AfterEach(func() {
		server.Close()
	})

	It("Should write messages", func() {
		conn := server.createConnection()
		under_test := NewWebsocketConnection(conn)

		under_test.WriteMessage("hi")

		received_messages, err := server.collectWebsocketMessagesSync()
		Expect(received_messages).To(ContainElement("hi"))
		Expect(err).To(BeNil())
	})

	Context("Client closing connections", func() {
		It("Should error when the client closes the connection", func() {
			test_connection := server.createConnection()
			under_test := NewWebsocketConnection(test_connection)

			test_connection.Close()

			Eventually(func() error { return under_test.WriteMessage("hi") }).ShouldNot(BeNil())
		})

		It("Should error when the client sends a close message", func() {
			test_connection := server.createConnection()
			under_test := NewWebsocketConnection(test_connection)

			test_connection.WriteMessage(websocket.CloseMessage, []byte{})

			Eventually(func() error { return under_test.WriteMessage("hi") }).ShouldNot(BeNil())
		})
	})

	Context("Server closing connections", func() {
		It("Should error when the server closes the connection", func() {
			test_connection := server.createConnection()
			under_test := NewWebsocketConnection(test_connection)

			server.closeConnections()

			Eventually(func() error { return under_test.WriteMessage("hi") }).ShouldNot(BeNil())
		})

		It("Should error when the server sends a close message", func() {
			test_connection := server.createConnection()
			under_test := NewWebsocketConnection(test_connection)

			server.writeCloseMessageToConnections()

			Eventually(func() error { return under_test.WriteMessage("hi") }).ShouldNot(BeNil())
		})
	})
})

type testServer struct {
	*httptest.Server
	writtenMessages  chan string
	connections_chan <-chan *websocket.Conn
	connections      []*websocket.Conn
	name             time.Time
	lock             sync.RWMutex
}

func (t *testServer) collectWebsocketMessagesSync() ([]string, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	messages := []string{}
	for {
		select {
		case message := <-t.writtenMessages:
			messages = append(messages, message)
		case <-ticker.C:
			if len(messages) == 0 {
				return nil, errors.New("Timeout receiving messages")
			}
			return messages, nil
		}
	}
}

func (t *testServer) createConnection() *websocket.Conn {
	dialer := &websocket.Dialer{}
	conn, _, _ := dialer.Dial(t.URL, nil)
	return conn
}

func (t *testServer) closeConnections() {
	t.lock.RLock()
	defer t.lock.RUnlock()
	for _, c := range t.connections {
		c.Close()
	}
}

func (t *testServer) writeCloseMessageToConnections() {
	t.lock.RLock()
	defer t.lock.RUnlock()
	for _, c := range t.connections {
		c.WriteMessage(websocket.CloseMessage, []byte{})
	}
}

func (t *testServer) writeMessageToConnections(message string) error {
	t.lock.RLock()
	defer t.lock.RUnlock()
	for _, c := range t.connections {
		err := c.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			return err
		}
	}
	return nil
}

func newServer() *testServer {
	messages := make(chan string, 10)
	connections_chan := make(chan *websocket.Conn, 10)
	connections := []*websocket.Conn{}
	s := httptest.NewServer(testHandler{messages, connections_chan})
	s.URL = makeWsProto(s.URL)
	test := &testServer{s, messages, connections_chan, connections, time.Now(), sync.RWMutex{}}

	go func() {
		for conn := range connections_chan {
			test.lock.Lock()
			test.connections = append(test.connections, conn)
			test.lock.Unlock()
		}
	}()

	return test
}

type testHandler struct {
	writtenMessages  chan string
	connections_chan chan *websocket.Conn
}

func (handler testHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	up := websocket.Upgrader{WriteBufferSize: 16}
	conn, _ := up.Upgrade(response, request, nil) // This closes the response writer
	handler.connections_chan <- conn

	for range time.Tick(10 * time.Millisecond) {
		conn.SetReadDeadline(time.Now().Add(time.Second))
		_, m, err := conn.ReadMessage()
		if string(m) != "" {
			handler.writtenMessages <- string(m)
		}
		if err != nil {
			return
		}
	}
}

func makeWsProto(s string) string {
	return "ws" + strings.TrimPrefix(s, "http")
}
