package main

import (
	"gopkg.in/gomail.v2"
)

func SendMail(smtp *Smtp, message StateChangeMessage) {
	// prepare
	m := gomail.NewMessage()
	m.SetHeader("From", smtp.From)
	m.SetHeader("To", smtp.To...)
	m.SetHeader("Subject", message.MailSubject())
	m.SetBody("text/html", message.MailBody())

	// connect
	d := gomail.NewDialer(smtp.Host, smtp.Port, smtp.Username, smtp.Password)

	// send & auto-close
	if err := d.DialAndSend(m); err != nil {
		// TODO : LOG ERROR
	}
}

