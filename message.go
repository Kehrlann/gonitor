package main

import (
	"fmt"
	"time"
	"text/template"
	"bytes"
)

// StateChangeMessage represents either a failure or a recovery
type StateChangeMessage struct {
	Resource Resource
	IsOk     bool
	Codes    []int
	Datetime time.Time
}

func (m *StateChangeMessage) String() string {
	return fmt.Sprintf("%v, %v, %v", m.Resource.Url, m.IsOk, m.Codes)
}

func (m *StateChangeMessage) MailSubject() string {
	messageType := "Failure"
	if m.IsOk {
		messageType = "Recovery"
	}
	return fmt.Sprintf("[Gonitor] %v for '%v'", messageType, m.Resource.Url)
}

func (m *StateChangeMessage) MailBody() string {
	error_recovery := "It seems an error occurred when polling <strong style=\"color:#d9534f\">{{.Resource.Url}}</strong>.<br>"
	if m.IsOk {
		error_recovery = "It seems <strong style=\"color:#5cb85c\">{{.Resource.Url}}</strong> has recovered.<br>"
	}
	template_string := "Hi !<br><br>" +
			"" +
			"This is an automated message from Gonitor.<br>" +
			error_recovery +
			"The following HTTP codes were received : {{.Codes}}.<br><br>" +
			"The config used is :" +
			"<ul>" +
			"<li>Interval : {{.Resource.IntervalInSeconds}} seconds</li>" +
			"<li>Number of tries : {{.Resource.NumberOfTries}}</li>" +
			"<li>Failure threshold : {{.Resource.FailureThreshold}}</li>" +
			"<li>Recovery threshold : {{.Resource.NumberOfTries}}</li>" +
			"<li>Timeout : {{.Resource.TimeoutInSeconds}} seconds</li>" +
			"</ul>"
	template, _ := template.New("body").Parse(template_string)
	var output bytes.Buffer
	template.Execute(&output, m)
	return output.String()
}

// ErrorMessage initializes a message with Time.Now() as the creation date and IsOk = false
func ErrorMessage(resource Resource, codes []int) *StateChangeMessage {
	m := &StateChangeMessage{resource, false, codes, time.Now()}
	return m
}

// NewMessage initializes a message with Time.Now() as the creation date and IsOk = true
func RecoveryMessage(resource Resource, codes []int) *StateChangeMessage {
	m := &StateChangeMessage{resource, true, codes, time.Now()}
	return m
}

// EmitMessages blah blah
func EmitMessages(messages <-chan *StateChangeMessage) {
	for m := range messages {
		fmt.Println(m)
	}
}
