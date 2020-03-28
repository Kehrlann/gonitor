package alert

import (
	"github.com/Kehrlann/gonitor/config"
	"github.com/Kehrlann/gonitor/monitor"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/Kehrlann/gonitor/websockets"
	"errors"
	"sync"
	"encoding/json"
)

var _ = Describe("websockets -> ", func() {

	// TODO : test errors !

	Describe("Registering / unregistering connections", func() {

		It("Should register a connection", func() {
			emitter := &WebsocketsEmitter{websocketConnections: make(map[uint]websockets.Connection)}
			conn := NewFake()

			emitter.registerConnection(conn)

			Expect(emitter.getConnections()).To(ContainElement(conn))
		})

		It("Should register multiple connetions", func() {
			emitter := &WebsocketsEmitter{websocketConnections: make(map[uint]websockets.Connection)}
			conn := NewFake()

			emitter.registerConnection(conn)
			emitter.registerConnection(conn)

			Expect(len(emitter.getConnections())).To(Equal(2))
		})

		It("Should unregister a connection", func() {
			emitter := &WebsocketsEmitter{websocketConnections: make(map[uint]websockets.Connection)}
			conn := NewFake()
			emitter.registerConnection(conn)

			emitter.unregisterConnection(0)

			Expect(emitter.getConnections()).To(BeEmpty())
		})

		It("Should not blow up when trying to unregister a non-indexed connection", func() {
			emitter := &WebsocketsEmitter{websocketConnections: make(map[uint]websockets.Connection)}
			conn := NewFake()
			emitter.registerConnection(conn)
			emitter.unregisterConnection(99)

			Expect(emitter.getConnections()).To(ContainElement(conn))
		})
	})

	Describe("getConnections", func() {
		It("Should copy the map", func() {
			emitter := &WebsocketsEmitter{websocketConnections: make(map[uint]websockets.Connection)}
			conn := NewFake()
			emitter.registerConnection(conn)

			connections := emitter.getConnections()
			delete(connections, 0)

			Expect(emitter.getConnections()).To(ContainElement(conn))
		})
	})

	Describe("New WebsocketEmitter", func() {
		var connections chan websockets.Connection
		var emitter *WebsocketsEmitter

		BeforeEach(func() {
			connections = make(chan websockets.Connection, 10)
			emitter = NewWebsocketEmitter(connections)
		})

		It("Should create", func() {
			Expect(emitter).ToNot(BeNil())
		})

		It("Should automatically register incoming connections", func() {
			conn := NewFake()
			connections <- conn

			Eventually(emitter.getConnections).Should(ContainElement(conn))
		})
	})

	Describe("Emit", func() {
		var connections chan websockets.Connection
		var emitter *WebsocketsEmitter
		var message *monitor.StateChangeMessage

		BeforeEach(func() {
			connections = make(chan websockets.Connection, 10)
			emitter = NewWebsocketEmitter(connections)
			message = monitor.RecoveryMessage(config.Resource{}, []int{})
		})

		It("Should not fail when there are no connections", func() {
			emitter.Emit(message)
		})

		It("Should write to all connections", func() {
			first_connection := NewFake()
			second_connection := NewFake()

			connections <- first_connection
			connections <- second_connection
			Eventually(emitter.getConnections).ShouldNot(BeEmpty())

			emitter.Emit(message)
			serialized_message, _ := json.Marshal(message)
			json_message := string(serialized_message)

			Eventually(first_connection.GetMessages).Should(ContainElement(json_message))
			Eventually(second_connection.GetMessages).Should(ContainElement(json_message))
		})

		It("Should unregister a erroring connections", func() {
			connections <- NewErrorFake()
			Eventually(emitter.getConnections).ShouldNot(BeEmpty()) // wait for the connection to be registered

			emitter.Emit(monitor.ErrorMessage(config.Resource{}, []int{}))
			Eventually(emitter.getConnections).Should(BeEmpty())
		})
	})
})

type FakeWebsocketConnection struct {
	messages    []string
	lock        sync.RWMutex
	shouldError bool
}

func NewFake() *FakeWebsocketConnection {
	return &FakeWebsocketConnection{[]string{}, sync.RWMutex{}, false}
}

func NewErrorFake() *FakeWebsocketConnection {
	return &FakeWebsocketConnection{[]string{}, sync.RWMutex{}, true}
}

func (conn *FakeWebsocketConnection) WriteMessage(message string) error {
	if conn.shouldError {
		return errors.New("error on write")
	}

	conn.lock.Lock()
	defer conn.lock.Unlock()
	conn.messages = append(conn.messages, message)
	return nil
}

func (conn *FakeWebsocketConnection) GetMessages() []string {
	conn.lock.RLock()
	defer conn.lock.RUnlock()
	return_array := []string{}
	for _, v := range conn.messages {
		return_array = append(return_array, v)
	}
	return return_array
}
