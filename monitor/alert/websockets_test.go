package alert

import (
	"github.com/kehrlann/gonitor/config"
	"github.com/kehrlann/gonitor/monitor"
	"github.com/gorilla/websocket"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("websockets -> ", func() {
	var message *monitor.StateChangeMessage
	// TODO : test me !
	// TODO : you need to read the message to be able to know whether the connection was closed or not ... test that
	// TODO : test timeout on client
	// TODO : test connection closed
	//			Caveat : apparently you have to read messages to be sure that the connection has been closed
	// 				conn.SetReadDeadline(time.Now().Add(time.Second))
	//				conn.ReadMessage()

	Describe("Registering / unregistering connections", func () {

		It("Should register a connection", func () {
			emitter := &WebsocketsEmitter{websocketConnections:make(map[uint]*websocket.Conn)}
			conn := &websocket.Conn{}

			emitter.RegisterConnection(conn)

			Expect(emitter.websocketConnections).To(ContainElement(conn))
		})

		It("Should register multiple connetions", func () {
			emitter := &WebsocketsEmitter{websocketConnections:make(map[uint]*websocket.Conn)}
			conn := &websocket.Conn{}

			emitter.RegisterConnection(conn)
			emitter.RegisterConnection(conn)

			Expect(len(emitter.websocketConnections)).To(Equal(2))
		})

		It("Should unregister a connection", func () {
			emitter := &WebsocketsEmitter{websocketConnections:make(map[uint]*websocket.Conn)}
			conn := &websocket.Conn{}
			emitter.RegisterConnection(conn)

			emitter.UnregisterConnection(0)

			Expect(emitter.websocketConnections).To(BeEmpty())
		})

		It("Should not blow up when trying to unregister a non-indexed connection", func () {
			emitter := &WebsocketsEmitter{websocketConnections:make(map[uint]*websocket.Conn)}
			conn := &websocket.Conn{}
			emitter.RegisterConnection(conn)
			emitter.UnregisterConnection(99)

			Expect(emitter.websocketConnections).To(ContainElement(conn))
		})
	})

	Describe("New WebsocketEmitter", func () {
		var connections chan *websocket.Conn
		var emitter *WebsocketsEmitter

		BeforeEach(func() {
			connections = make(chan *websocket.Conn, 10)
			emitter = NewWebsocketEmitter(connections)
		})

		It("Should create", func() {
			Expect(emitter).ToNot(BeNil())
		})

		It("Should manually register one connection", func() {
			conn := &websocket.Conn{}
			emitter.RegisterConnection(conn)

			Expect(emitter.websocketConnections).To(ContainElement(conn))
		})

		It("Should manually register multiple connections", func() {
			emitter.RegisterConnection(&websocket.Conn{})
			emitter.RegisterConnection(&websocket.Conn{})
			emitter.RegisterConnection(&websocket.Conn{})

			Expect(len(emitter.websocketConnections)).To(BeNumerically(">", 1))
		})

		It("Should manually unregister an existing connection", func() {
			conn := &websocket.Conn{}
			emitter.RegisterConnection(conn)
			emitter.UnregisterConnection(emitter.currentIndex - 1)

			Expect(emitter.websocketConnections).ToNot(ContainElement(conn))
		})

		It("Should not fail when trying to unregister a non existing connection", func() {
			conn := &websocket.Conn{}
			emitter.RegisterConnection(conn)
			emitter.UnregisterConnection(999)

			Expect(emitter.websocketConnections).To(ContainElement(conn))
		})

		It("Should automatically register incoming connections", func() {
			conn := &websocket.Conn{}
			connections <- conn

			Eventually(emitter.websocketConnections).Should(ContainElement(conn))
		})
	})

	BeforeEach(func() {
		res := config.Resource{"http://test.com", 60, 2, 10, 3, "" }
		message = monitor.RecoveryMessage(res, []int{1, 2, 3})
	})

	Context("Tech tests ", func() {
		It("is fun", func() {
		})
	})
})
