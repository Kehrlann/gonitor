package emit

import (
	log "github.com/Sirupsen/logrus"
)

// LogEmitter logs the messages in the console
type LogEmitter struct {
}

// Emit sends a StateChangeMessage via e-mail
func (emitter *LogEmitter) Emit(message *StateChangeMessage) {
	log.Info(message)
}
