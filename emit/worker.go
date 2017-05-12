package emit

import (
	"github.com/kehrlann/gonitor/config"
)

func EmitMessages(messages <-chan *StateChangeMessage, configuration *config.Configuration) {
	emitters := []Emitter{
		&MailEmitter{&configuration.Smtp},
		&CommandEmitter{configuration.GlobalCommand},
		&LogEmitter{},
	}
	emitMessages(emitters, messages)
}

func emitMessages(emitters []Emitter, messages <-chan *StateChangeMessage){
	for m := range messages {
		for _, e := range emitters {
			go e.Emit(m)
		}
	}
}

