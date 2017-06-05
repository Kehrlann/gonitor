package alert

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kehrlann/gonitor/config"

	"gopkg.in/gomail.v2"
	"github.com/kehrlann/gonitor/monitor"
)

// Mail emitter emits messages via e-mail
type mailEmitter struct {
	smtp *config.Smtp
}

// mailer describes an entity able to send e-mails. Mostly used for testing purposes
type mailer interface {
	DialAndSend(mail ... *gomail.Message) error
}

// Emit sends a StateChangeMessage via e-mail
func (emitter *mailEmitter) Emit(message *monitor.StateChangeMessage) {
	mailer := gomail.NewDialer(emitter.smtp.Host, emitter.smtp.Port, emitter.smtp.Username, emitter.smtp.Password)
	sendMail(mailer, emitter.smtp, message)
}

func sendMail(mailer mailer, smtp *config.Smtp, message *monitor.StateChangeMessage) {
	if !smtp.IsValid() {
		return
	}

	// prepare
	m := gomail.NewMessage()
	m.SetHeader("From", smtp.FormatFromHeader())
	m.SetHeader("To", smtp.To...)
	m.SetHeader("Subject", message.MailSubject())
	m.SetBody("text/html", message.MailBody())

	// send & auto-close
	if err := mailer.DialAndSend(m); err != nil {
		log.Errorf("Error sending e-mail alert : `%v`", err)
	}
}
