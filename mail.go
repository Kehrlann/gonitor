package main

import (
	"gopkg.in/gomail.v2"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

func SendMail(smtp *Smtp, message *StateChangeMessage) {
	if !smtp.IsValid() {
		return
	}

	// prepare
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%v <%v>", smtp.FromName, smtp.FromAddress))
	m.SetHeader("To", smtp.To...)
	m.SetHeader("Subject", message.MailSubject())
	m.SetBody("text/html", message.MailBody())

	// connect
	d := gomail.NewDialer(smtp.Host, smtp.Port, smtp.Username, smtp.Password)

	// send & auto-close
	if err := d.DialAndSend(m); err != nil {
		log.Errorf("Error sending e-mail alert : `%v`", err)
	}
}
