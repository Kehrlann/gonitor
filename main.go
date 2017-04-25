package main

import (
	log "github.com/Sirupsen/logrus"
	"flag"
	"fmt"
)

func main() {
	// TODO : rework, test
	flag.Usage = func() {
		fmt.Println()
		fmt.Println("Usage: gonitor [-config path_to_config]")
		fmt.Println()
		flag.PrintDefaults()
	}
	path := flag.String("config", "./gonitor.config.json", "Path to the config file")
	flag.Parse()



	log.SetLevel(log.DebugLevel)
	log.Info("Starting Gonitor ...")
	config, err := LoadConfig(*path)

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
