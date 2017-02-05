package main

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func consumeMessages(channel <-chan *Message) {
	for _ = range channel {

	}
}

var _ = Describe("Analyze", func() {

	Describe("Computing state", func() {})

	Describe("Analyzing return codes", func() {
		Context("Given only successes", func() {
			It("Should not emit messages", func() {
				responseCodes := make(chan int)
				messages := make(chan *Message)

				go Consistently(messages, 1*time.Second).ShouldNot(Receive())

				go Analyze("fake_url", responseCodes, messages)

				for i := 1; i <= 20; i++ {
					responseCodes <- 200
				}
			})
		})

		Context("Given only errors", func() {

			It("Should not emit a message before threshold reached", func() {
				responseCodes := make(chan int)
				messages := make(chan *Message)

				go Consistently(messages, 1*time.Second).ShouldNot(Receive())

				go Analyze("fake_url", responseCodes, messages)

				for i := 1; i <= 2; i++ {
					responseCodes <- 404
				}
			})

			It("Should emit a message when threshold reached", func() {
				responseCodes := make(chan int)
				messages := make(chan *Message)

				go Analyze("should emit", responseCodes, messages)
				go consumeMessages(messages)

				for i := 1; i <= 4; i++ {
					responseCodes <- 404
				}

				var receivedMessage *Message
				Eventually(messages, 1*time.Second).Should(Receive(&receivedMessage))
				Expect(receivedMessage.String()).To(ContainSubstring("false"))
			})
		})

		Context("Given successes and errors", func() {
			It("Should emit a message", func() {
				responseCodes := make(chan int)
				messages := make(chan *Message)

				var receivedMessage *Message
				go Eventually(messages).Should(Receive(&receivedMessage))

				go Analyze("fake_url", responseCodes, messages)

				for i := 1; i <= 10; i++ {
					responseCodes <- 200
				}
				for i := 1; i <= 10; i++ {
					responseCodes <- 404
				}

				Expect(receivedMessage.String()).To(ContainSubstring("false"))
			})
		})

	})

})
