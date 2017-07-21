package alert

import (
	"github.com/kehrlann/gonitor/config"
	testlog "github.com/sirupsen/logrus/hooks/test"
	log "github.com/sirupsen/logrus"
	
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gopkg.in/gomail.v2"
	"errors"
	"github.com/kehrlann/gonitor/monitor"
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
	var message *monitor.StateChangeMessage
	var mailer *FakeMailer

	BeforeEach(func() {
		res := config.Resource{"http://test.com", 60, 2, 10, 3, "" }
		// TODO : use pointers to resources ?
		message = monitor.RecoveryMessage(res, []int{1, 2, 3})
		mailer = &FakeMailer{}
	})

	Context("When SMTP isn't valid -> ", func() {
		It("Shouldn't try to send an e-mail", func() {
			smtp := &config.Smtp{}
			sendMail(mailer, smtp, message)
			Expect(mailer.messages).To(BeEmpty())
		})
	})

	Context("When SMTP is valid -> ", func() {
		smtp := &config.Smtp{"mail.example.com", 25, "user", "password", "user@example.com", "User", []string{"recipient@example.com" }}

		It("Should send an e-mail", func() {
			sendMail(mailer, smtp, message)

			Expect(len(mailer.messages)).To(Equal(1))
			// mmmmh can't verify the struct ¬.¬
		})

		It("Should log when something goes wrong", func() {
			hook := testlog.NewGlobal()
			log.SetLevel(log.ErrorLevel)
			log.SetOutput(&nilWriter{})
			mailer := &FakeMailer{
				make([]*gomail.Message, 0),
				func(m *gomail.Message) error { return errors.New("Woops") },
			}

			sendMail(mailer, smtp, message)

			Eventually(func() []*log.Entry { return hook.Entries }).ShouldNot(BeEmpty())
			Expect(hook.LastEntry().Level).To(Equal(log.ErrorLevel))
		})
	})
})

type nilWriter struct {}

func (w *nilWriter) Write(p []byte) (int, error) {
	return 0, nil
}
