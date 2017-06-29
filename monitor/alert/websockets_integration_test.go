package alert

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http/httptest"
	"strings"
	"net/http"
	"github.com/gorilla/websocket"
	"time"
	"errors"
	"github.com/kehrlann/gonitor/monitor"
	"github.com/kehrlann/gonitor/config"
)

var _ = Describe("Websockets integrations tests -> ", func() {

	var server testServer
	message := monitor.RecoveryMessage(config.Resource{}, []int{})
	var emitter *WebsocketsEmitter
	var conn chan *websocket.Conn

	BeforeEach(func() {
		server = newServer()
		conn = make(chan *websocket.Conn, 10)
		emitter = NewWebsocketEmitter(conn)
	})

	AfterEach(func() {
		server.Close()
	})

	It("It should have a test server", func() {
		conn := server.createConnection()
		conn.WriteMessage(websocket.TextMessage, []byte("hi"))
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
		conn <- server.createConnection()
		Eventually(emitter.getConnections).ShouldNot(BeEmpty()) // Wait for connections to register in the emitter

		emitter.Emit(message)
		messages, err := server.collectWebsocketMessagesSync()

		Expect(messages).ToNot(BeEmpty())
		Expect(len(messages)).To(Equal(1))
		Expect(err).To(BeNil())
	})

	It("Should send multiple messages when there on one connection", func() {
		conn <- server.createConnection()
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
		conn <- server.createConnection()
		conn <- server.createConnection()
		Eventually(countConnections).Should(Equal(2)) // Wait for connections to register in the emitter

		emitter.Emit(message)
		messages, err := server.collectWebsocketMessagesSync()

		Expect(messages).ToNot(BeEmpty())
		Expect(len(messages)).To(Equal(2))
		Expect(err).To(BeNil())
	})

	It("Should unregister a closed connection", func() {
		test_connection := server.createConnection()
		conn <- test_connection
		Eventually(countConnections).Should(Equal(1)) // Wait for connections to register in the emitter

		test_connection.Close()
		Consistently(countConnections).Should(Equal(1))

		emitter.Emit(message)
		Eventually(countConnections).Should(Equal(0))
	})

	It("Should unregister a connection when receiving a `Close` control message", func() {
		test_connection := server.createConnection()
		conn <- test_connection
		Eventually(countConnections).Should(Equal(1)) // Wait for connections to register in the emitter

		test_connection.WriteMessage(websocket.CloseMessage, []byte{})
		Consistently(countConnections).Should(Equal(1))

		emitter.Emit(message)
		Eventually(countConnections).Should(Equal(0))
	})
})

type testServer struct {
	*httptest.Server
	writtenMessages chan string
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

func newServer() testServer {
	messages := make(chan string, 10)
	s := httptest.NewServer(testHandler{messages})
	s.URL = makeWsProto(s.URL)
	test := testServer{s, messages}
	return test
}

type testHandler struct {
	writtenMessages chan string
}

func (handler testHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	up := websocket.Upgrader{WriteBufferSize: 1024}
	conn, _ := up.Upgrade(response, request, nil) // This closes the response writer

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
