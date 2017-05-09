package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gopkg.in/gomail.v2"
	"errors"
	testlog "github.com/Sirupsen/logrus/hooks/test"
	log "github.com/Sirupsen/logrus"
)

type FakeMailer struct {
	messages      []*gomail.Message
	handleMessage func(message *gomail.Message) error
}

func (mailer *FakeMailer) DialAndSend(messages ... *gomail.Message) error {
	handler := func(m *gomail.Message) error { return nil }
	if mailer.handleMessage != nil {
		handler = mailer.handleMessage
	}

	for _, m := range messages {
		mailer.messages = append(mailer.messages, m)
		if err := handler(m); err != nil {
			return err
		}
	}

	return nil
}

var _ = Describe("sendMail -> ", func() {

	var message *StateChangeMessage
	var mailer *FakeMailer

	BeforeEach(func() {
		res := Resource{"http://test.com", 60, 2, 10, 3, "" }
		// TODO : use pointers to resources ?
		message = RecoveryMessage(res, []int{1, 2, 3})
		mailer = &FakeMailer{}
	})

	Context("When SMTP isn't valid -> ", func() {
		It("Shouldn't try to send an e-mail", func() {
			smtp := &Smtp{}
			sendMail(mailer, smtp, message)
			Expect(mailer.messages).To(BeEmpty())
		})
	})

	Context("When SMTP is valid -> ", func() {
		smtp := &Smtp{"mail.example.com", 25, "user", "password", "user@example.com", "User", []string{"recipient@example.com" }}

		It("Should send an e-mail", func() {
			sendMail(mailer, smtp, message)

			Expect(len(mailer.messages)).To(Equal(1))
			// mmmmh can't verify the struct ¬.¬
		})

		It("Should log when something goes wrong", func() {
			hook := testlog.NewGlobal()
			log.SetLevel(log.ErrorLevel)

			mailer = &FakeMailer{
				make([]*gomail.Message, 0),
				func(m *gomail.Message) error { return errors.New("Woops") },
			}

			sendMail(mailer, smtp, message)

			Expect(hook.Entries).ToNot(BeEmpty())
			Expect(hook.LastEntry().Level).To(Equal(log.ErrorLevel))
		})
	})

})
