package alert

import "github.com/Kehrlann/gonitor/monitor"

type emitter interface {
	Emit(message *monitor.StateChangeMessage)
}
