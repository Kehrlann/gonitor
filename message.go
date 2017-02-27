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

// ErrorMessage initializes a message with Time.Now() as the creation date and isOk = false
func ErrorMessage(url string, codes []int) *StateChangeMessage {
	m := &StateChangeMessage{url, false, codes, time.Now()}
	return m
}

// NewMessage initializes a message with Time.Now() as the creation date and isOk = true
func RecoveryMessage(url string, codes []int) *StateChangeMessage {
	m := &StateChangeMessage{url, true, codes, time.Now()}
	return m
}

// EmitMessages blah blah
func EmitMessages(messages <-chan *StateChangeMessage) {
	for m := range messages {
		fmt.Println(m)
	}
}
