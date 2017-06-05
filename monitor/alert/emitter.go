package alert

import "github.com/kehrlann/gonitor/monitor"

type emitter interface {
	Emit(message *monitor.StateChangeMessage)
}
