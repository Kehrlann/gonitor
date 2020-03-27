package alert

import (
	"github.com/Kehrlann/gonitor/config"
	"github.com/Kehrlann/gonitor/monitor"
	"github.com/Kehrlann/gonitor/websockets"
)

// EmitMessages emits all messages to STDOUT, sends alerts via e-mail and executes the configured commants
func EmitMessages(messages <-chan *monitor.StateChangeMessage,
	websocketsConnections <- chan websockets.Connection,
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
