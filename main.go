package main

import (
	"log"
)

func main() {
	config, err := LoadConfig("./config.json")

	if err != nil {
		log.Fatal(err)
	}

	messages := make(chan *StateChangeMessage)
	for _, resource := range config.Resources {
		go Run(resource, messages)
	}
	EmitMessages(messages)
}
