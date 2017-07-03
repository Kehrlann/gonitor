package alert

import (
	"net/http/httptest"
	"time"
	"strings"

	"github.com/kehrlann/gonitor/config"
	"github.com/kehrlann/gonitor/monitor"
	"github.com/kehrlann/gonitor/websockets"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/gorilla/websocket"
	"errors"
	"net/http"
)

var _ = Describe("Websockets integrations tests -> ", func() {

	var server *testServer
	message := monitor.RecoveryMessage(config.Resource{}, []int{})
	var emitter *WebsocketsEmitter
	var conn chan *websockets.Connection

	BeforeEach(func() {
		server = newServer()
		conn = make(chan *websockets.Connection, 10)
		emitter = NewWebsocketEmitter(conn)
	})

	AfterEach(func() {
		server.Close()
	})

	It("It should have a test server", func() {
		conn := server.createConnection()
		conn.WriteMessage("hi")
		received_messages, err := server.collectWebsocketMessagesSync()

		Expect(received_messages).To(ContainElement("hi"))
		Expect(err).To(BeNil())
	})

	It("Should work when there are no connections", func() {
		emitter.Emit(message)
		messages, err := server.collectWebsocketMessagesSync()

		Expect(messages).To(BeNil())
		Expect(err).ToNot(BeNil())
	})

	It("Should send messages when there is one connection", func() {
		pushConnection(conn, server.createConnection())
		Eventually(emitter.getConnections).ShouldNot(BeEmpty()) // Wait for connections to register in the emitter

		emitter.Emit(message)
		messages, err := server.collectWebsocketMessagesSync()

		Expect(messages).ToNot(BeEmpty())
		Expect(len(messages)).To(Equal(1))
		Expect(err).To(BeNil())
	})

	It("Should send multiple messages when there on one connection", func() {
		pushConnection(conn, server.createConnection())
		Eventually(emitter.getConnections).ShouldNot(BeEmpty()) // Wait for connections to register in the emitter

		emitter.Emit(message)
		emitter.Emit(message)
		messages, err := server.collectWebsocketMessagesSync()

		Expect(messages).ToNot(BeEmpty())
		Expect(len(messages)).To(Equal(2))
		Expect(err).To(BeNil())
	})

	countConnections := func() int { return len(emitter.getConnections()) }
	It("Should send messages when there are multiple connections", func() {
		pushConnection(conn, server.createConnection())
		pushConnection(conn, server.createConnection())
		Eventually(countConnections).Should(Equal(2)) // Wait for connections to register in the emitter

		emitter.Emit(message)
		messages, err := server.collectWebsocketMessagesSync()

		Expect(messages).ToNot(BeEmpty())
		Expect(len(messages)).To(Equal(2))
		Expect(err).To(BeNil())
	})

	It("Should not blow up on concurrent writes", func() {
		test_connection := server.createConnection()
		conn <- &test_connection
		Eventually(countConnections).Should(Equal(1)) // Wait for connections to register in the emitter

		go emitter.Emit(message)
		go emitter.Emit(message)
		messages, err := server.collectWebsocketMessagesSync()

		Expect(messages).ToNot(BeEmpty())
		Expect(len(messages)).To(Equal(2))
		Expect(err).To(BeNil())
	})

	// TODO : doesn't detect close when we don't write !
	Context("Server closing connections", func () {
		It("Should unregister a closed connection", func() {
			test_connection := server.createConnection()
			conn <- &test_connection
			Eventually(countConnections).Should(Equal(1)) // Wait for connections to register in the emitter

			server.closeConnections()
			Eventually(countConnections).Should(Equal(0))
		})

		It("Should unregister a connection when receiving a `Close` control message", func() {
			test_connection := server.createConnection()
			conn <- &test_connection
			Eventually(countConnections).Should(Equal(1)) // Wait for connections to register in the emitter

			server.writeCloseMessageToConnections()
			Eventually(countConnections).Should(Equal(0))
		})
	})
})

type testServer struct {
	*httptest.Server
	writtenMessages chan string
	connections_chan <-chan *websocket.Conn
	connections []*websocket.Conn
	 name time.Time
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

func (t *testServer) createConnection() websockets.Connection {
	dialer := &websocket.Dialer{}
	conn, _, _ := dialer.Dial(t.URL, nil)

	return websockets.NewWebsocketConnection(conn)
}

func (t *testServer) closeConnections() {
	for _, c := range t.connections {
		c.Close()
	}
}

func (t *testServer) writeCloseMessageToConnections() {
	for _, c := range t.connections {
		c.WriteMessage(websocket.CloseMessage, []byte{})
	}
}

func newServer() *testServer {
	messages := make(chan string, 10)
	connections_chan := make(chan *websocket.Conn, 10)
	connections := []*websocket.Conn{}
	s := httptest.NewServer(testHandler{messages, connections_chan})
	s.URL = makeWsProto(s.URL)
	test := &testServer{s, messages, connections_chan, connections, time.Now()}

	go func() {
		for conn := range connections_chan {
			test.connections = append(test.connections, conn)
		}
	}()

	return test
}

type testHandler struct {
	writtenMessages chan string
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

func pushConnection(connections chan *websockets.Connection, connection websockets.Connection) {
	connections <- &connection
}