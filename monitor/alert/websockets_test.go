package alert

import (
	"github.com/kehrlann/gonitor/config"
	"github.com/kehrlann/gonitor/monitor"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/kehrlann/gonitor/websockets"
	"errors"
	"sync"
)

var _ = Describe("websockets -> ", func() {

	// TODO : test errors !

	Describe("Registering / unregistering connections", func () {

		It("Should register a connection", func () {
			emitter := &WebsocketsEmitter{websocketConnections:make(map[uint]*websockets.Connection)}
			conn := NewFake()

			emitter.registerConnection(&conn)

			Expect(emitter.getConnections()).To(ContainElement(&conn))
		})

		It("Should register multiple connetions", func () {
			emitter := &WebsocketsEmitter{websocketConnections:make(map[uint]*websockets.Connection)}
			conn := NewFake()

			emitter.registerConnection(&conn)
			emitter.registerConnection(&conn)

			Expect(len(emitter.getConnections())).To(Equal(2))
		})

		It("Should unregister a connection", func () {
			emitter := &WebsocketsEmitter{websocketConnections:make(map[uint]*websockets.Connection)}
			conn := NewFake()
			emitter.registerConnection(&conn)

			emitter.unregisterConnection(0)

			Expect(emitter.getConnections()).To(BeEmpty())
		})

		It("Should not blow up when trying to unregister a non-indexed connection", func () {
			emitter := &WebsocketsEmitter{websocketConnections:make(map[uint]*websockets.Connection)}
			conn := NewFake()
			emitter.registerConnection(&conn)
			emitter.unregisterConnection(99)

			Expect(emitter.getConnections()).To(ContainElement(&conn))
		})
	})

	Describe("getConnections", func () {
		It("Should copy the map", func () {
			emitter := &WebsocketsEmitter{websocketConnections:make(map[uint]*websockets.Connection)}
			conn := NewFake()
			emitter.registerConnection(&conn)

			connections := emitter.getConnections()
			delete(connections, 0)

			Expect(emitter.getConnections()).To(ContainElement(&conn))
		})
	})

	Describe("New WebsocketEmitter", func () {
		var connections chan *websockets.Connection
		var emitter *WebsocketsEmitter

		BeforeEach(func() {
			connections = make(chan *websockets.Connection, 10)
			emitter = NewWebsocketEmitter(connections)
		})

		It("Should create", func() {
			Expect(emitter).ToNot(BeNil())
		})

		It("Should automatically register incoming connections", func() {
			conn := NewFake()
			connections <- &conn

			Eventually(emitter.getConnections).Should(ContainElement(&conn))
		})
	})

	Describe("Emit", func() {
		var connections chan *websockets.Connection
		var emitter *WebsocketsEmitter
		var message *monitor.StateChangeMessage

		BeforeEach(func() {
			connections = make(chan *websockets.Connection, 10)
			emitter = NewWebsocketEmitter(connections)
			message = monitor.RecoveryMessage(config.Resource{}, []int {})
		})

		It("Should not fail when there are no connections", func () {
			emitter.Emit(message)
		})

		It("Should write to all connections", func() {
			first_connection := NewFake()
			second_connection := NewFake()

			connections <- &first_connection
			connections <- &second_connection
		})
	})
})

type FakeWebsocketConnection struct {
	messages []string
	lock sync.RWMutex
}

func NewFake() websockets.Connection {
	return &FakeWebsocketConnection{[]string{}, sync.RWMutex{}}
}

func (conn *FakeWebsocketConnection) WriteMessage(message string) error {
	if message == "error" {
		return errors.New("error on write")
	}
	
	conn.lock.Lock()
	defer conn.lock.Unlock()
	conn.messages = append(conn.messages, message)
	return nil
}
