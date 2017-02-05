package main

import (
	"fmt"
	"time"
)

// Message represents either a failure or a recovery
type Message struct {
	url      string
	isOk     bool
	codes    []int
	datetime time.Time
}

func (m *Message) String() string {
	return fmt.Sprintf("%v, %v, %v", m.url, m.isOk, m.codes)
}

// NewMessage initializes a message with Time.Now() as the creation date
func NewMessage(url string, isOk bool, codes []int) *Message {
	m := &Message{url, isOk, codes, time.Now()}
	return m
}

// EmitMessages blah blah
func EmitMessages(messages <-chan *Message) {
	for m := range messages {
		fmt.Println(m)
	}
}
