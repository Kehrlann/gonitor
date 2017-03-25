package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	config, err := LoadConfig("./config.json")

	if err != nil {
		log.Fatal(err)
	}

	// TODO : use config to start Goroutines
	fmt.Println(config)

	messages := make(chan *StateChangeMessage)
	for _, resource := range config.Resources {
		go Run(resource, 2*time.Second, messages)
	}
	EmitMessages(messages)
}
