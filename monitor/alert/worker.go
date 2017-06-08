package alert

import (
	"github.com/kehrlann/gonitor/config"
	"github.com/kehrlann/gonitor/monitor"
	"github.com/gorilla/websocket"
)

// EmitMessages emits all messages to STDOUT, sends alerts via e-mail and executes the configured commants
func EmitMessages(messages <-chan *monitor.StateChangeMessage,
	websocketsConnections <- chan *websocket.Conn,
	configuration *config.Configuration) {
	emitters := []emitter{
		&mailEmitter{&configuration.Smtp},
		&commandEmitter{configuration.GlobalCommand},
		&logEmitter{},
		NewWebsocketEmitter(websocketsConnections),
	}
	emitMessages(emitters, messages)
}

func emitMessages(emitters []emitter, messages <-chan *monitor.StateChangeMessage) {
	for m := range messages {
		for _, e := range emitters {
			go e.Emit(m)
		}
	}
}
