package main

import "time"

func main() {
	messages := make(chan *Message)
	go Run("https://ratp.garnier.wf", 2*time.Second, messages)
	go Run("http://www.google.com", 2*time.Second, messages)
	EmitMessages(messages)
}
