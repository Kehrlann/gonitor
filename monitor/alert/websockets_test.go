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

			emitter.registerConnection(conn)

			Expect(emitter.getConnections()).To(ContainElement(conn))
		})

		It("Should register multiple connetions", func () {
			emitter := &WebsocketsEmitter{websocketConnections:make(map[uint]*websocket.Conn)}
			conn := &websocket.Conn{}

			emitter.registerConnection(conn)
			emitter.registerConnection(conn)

			Expect(len(emitter.getConnections())).To(Equal(2))
		})

		It("Should unregister a connection", func () {
			emitter := &WebsocketsEmitter{websocketConnections:make(map[uint]*websocket.Conn)}
			conn := &websocket.Conn{}
			emitter.registerConnection(conn)

			emitter.unregisterConnection(0)

			Expect(emitter.getConnections()).To(BeEmpty())
		})

		It("Should not blow up when trying to unregister a non-indexed connection", func () {
			emitter := &WebsocketsEmitter{websocketConnections:make(map[uint]*websocket.Conn)}
			conn := &websocket.Conn{}
			emitter.registerConnection(conn)
			emitter.unregisterConnection(99)

			Expect(emitter.getConnections()).To(ContainElement(conn))
		})
	})

	Describe("getConnections", func () {
		It("Should copy the map", func () {
			emitter := &WebsocketsEmitter{websocketConnections:make(map[uint]*websocket.Conn)}
			conn := &websocket.Conn{}
			emitter.registerConnection(conn)

			connections := emitter.getConnections()
			delete(connections, 0)

			Expect(emitter.getConnections()).To(ContainElement(conn))
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

		It("Should automatically register incoming connections", func() {
			conn := &websocket.Conn{}
			connections <- conn

			Eventually(emitter.getConnections).Should(ContainElement(conn))
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
