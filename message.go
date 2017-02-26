package main

import (
	"fmt"
	"time"
)

// StateChangeMessage represents either a failure or a recovery
type StateChangeMessage struct {
	url      string
	isOk     bool
	codes    []int
	datetime time.Time
}

func (m *StateChangeMessage) String() string {
	return fmt.Sprintf("%v, %v, %v", m.url, m.isOk, m.codes)
}

// NewMessage initializes a message with Time.Now() as the creation date
func NewMessage(url string, isOk bool, codes []int) *StateChangeMessage {
	m := &StateChangeMessage{url, isOk, codes, time.Now()}
	return m
}

// EmitMessages blah blah
func EmitMessages(messages <-chan *StateChangeMessage) {
	for m := range messages {
		fmt.Println(m)
	}
}
