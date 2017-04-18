package main

import (
	"log"
	"fmt"
)

func main() {
	config, err := LoadConfig("./gonitor.config.json")

	if err != nil {
		log.Fatal(err)
	}

	messages := make(chan *StateChangeMessage)
	for _, resource := range config.Resources {
		go Run(resource, messages)
	}
	EmitMessages(messages, &config.Smtp)
}


// EmitMessages blah blah
func EmitMessages(messages <-chan *StateChangeMessage, smtp *Smtp) {
	for m := range messages {
		fmt.Println(m)
		go SendMail(smtp, m)
	}
}
