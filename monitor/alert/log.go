package alert

import (
	log "github.com/sirupsen/logrus"
	"github.com/kehrlann/gonitor/monitor"
)

// logEmitter logs the messages in the console
type logEmitter struct {
}

// Emit sends a StateChangeMessage via e-mail
func (emitter *logEmitter) Emit(message *monitor.StateChangeMessage) {
	log.Info(message)
}
