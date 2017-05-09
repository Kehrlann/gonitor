package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"

	"gopkg.in/gomail.v2"
)


type Mailer interface {
	DialAndSend(mail... *gomail.Message) error
}

// SendMail sends a StateChangeMessage via e-mail
func SendMail(smtp *Smtp, message *StateChangeMessage) {
	mailer := gomail.NewDialer(smtp.Host, smtp.Port, smtp.Username, smtp.Password)
	sendMail(mailer, smtp, message)
}

func sendMail(mailer Mailer, smtp *Smtp, message *StateChangeMessage) {
	if !smtp.IsValid() {
		return
	}

	// prepare
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%v <%v>", smtp.FromName, smtp.FromAddress))
	m.SetHeader("To", smtp.To...)
	m.SetHeader("Subject", message.MailSubject())
	m.SetBody("text/html", message.MailBody())

	// send & auto-close
	if err := mailer.DialAndSend(m); err != nil {
		log.Errorf("Error sending e-mail alert : `%v`", err)
	}
}
