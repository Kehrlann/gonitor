package main

import (
	"gopkg.in/gomail.v2"
)

func SendMail(smtp *Smtp, message *StateChangeMessage) {
	// TODO : HERE
	// TODO : HERE
	m := gomail.NewMessage()
	m.SetHeader("From", smtp.From)
	m.SetHeader("To", smtp.To...)
	m.SetHeader("Subject", "Gonitor : Error for ")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")

	d := gomail.NewDialer(smtp.Host, smtp.Port, smtp.Username, smtp.Password)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		// TODO : LOG ERROR
	}
}

