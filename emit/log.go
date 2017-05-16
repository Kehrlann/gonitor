package emit

import (
	log "github.com/Sirupsen/logrus"
)

// logEmitter logs the messages in the console
type logEmitter struct {
}

// Emit sends a StateChangeMessage via e-mail
func (emitter *logEmitter) Emit(message *StateChangeMessage) {
	log.Info(message)
}
