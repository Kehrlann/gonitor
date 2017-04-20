package main

import (
	log "github.com/Sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Info("Starting Gonitor ...")
	config, err := LoadConfig("./gonitor.config.json")

	if err != nil {
		log.Fatalf("Error loading config : `%v`", err)
	}

	log.Info("Starting monitoring ...")
	messages := make(chan *StateChangeMessage)
	for _, resource := range config.Resources {
		go Run(resource, messages)
	}
	EmitMessages(messages, &config.Smtp)
}

// EmitMessages blah blah
func EmitMessages(messages <-chan *StateChangeMessage, smtp *Smtp) {
	for m := range messages {
		log.Info(m)
		go SendMail(smtp, m)
	}
}
