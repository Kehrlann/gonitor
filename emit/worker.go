package emit

import (
	"github.com/kehrlann/gonitor/config"
	log "github.com/Sirupsen/logrus"
)

func EmitMessages(messages <-chan *StateChangeMessage, configuration *config.Configuration) {
	for m := range messages {
		log.Info(m)
		go SendMail(&configuration.Smtp, m)
		go ExecCommand(m, configuration.GlobalCommand)
	}
}

