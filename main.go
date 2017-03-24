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
	go Run("https://ratp.garnier.wf", 2*time.Second, messages)
	go Run("http://www.google.com", 2*time.Second, messages)
	EmitMessages(messages)
}
