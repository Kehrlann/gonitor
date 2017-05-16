package monitor

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/kehrlann/gonitor/config"
	"github.com/kehrlann/gonitor/emit"
)

var _ = Describe("analyze", func() {

	Describe("Computing state", func() {

		Context("Given an empty array", func() {
			codesToAnalyze := []int{0}

			It("Should not be a failure", func() {
				failure, _ := computeState(codesToAnalyze, 3, 10)
				Expect(failure).To(BeFalse())
			})

			It("Should not allow recovery", func() {
				_, canRecover := computeState(codesToAnalyze, 3, 10)
				Expect(canRecover).To(BeFalse())
			})
		})

		Context("Given all successes", func() {
			codesToAnalyze := []int{200, 200, 200, 200, 200, 200, 200, 200, 200, 200}

			It("Should not be a failure", func() {
				failure, _ := computeState(codesToAnalyze, 3, 10)
				Expect(failure).To(BeFalse())
			})

			It("Should allow recovery", func() {
				_, canRecover := computeState(codesToAnalyze, 3, len(codesToAnalyze))
				Expect(canRecover).To(BeTrue())
			})
		})

		Context("Given one failure", func() {

			codesToAnalyze := []int{200, 200, 200, 200, 200, 200, 200, 200, 200, 0}
			It("Should not be a failure", func() {
				failure, _ := computeState(codesToAnalyze, 3, 10)
				Expect(failure).To(BeFalse())
			})

			It("Should not allow recovery", func() {
				_, canRecover := computeState(codesToAnalyze, 3, 10)
				Expect(canRecover).To(BeFalse())
			})
		})

		Context("Given all failures", func() {
			codesToAnalyze := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

			It("Should be a failure", func() {
				failure, _ := computeState(codesToAnalyze, 3, 10)
				Expect(failure).To(BeTrue())
			})

			It("Should not allow recovery", func() {
				_, canRecover := computeState(codesToAnalyze, 3, 10)
				Expect(canRecover).To(BeFalse())
			})
		})

	})

	Describe("Receiving HTTP response Codes", func() {

		resource := config.Resource{"Url", 2, 2, 3, 2, ""}

		Context("When not polling", func() {
			It("Should not emit messages", func() {
				messages := make(chan *emit.StateChangeMessage)
				codes := make(chan int)
				go analyze(resource, codes, messages)
				Consistently(messages).ShouldNot(Receive())
			})
		})

		Context("When getting only successes", func() {

			emitHttpOk := func(codes chan<- int) {
				for range time.Tick(2 * time.Millisecond) {
					// TODO : emit random Codes between 200 and 300 (HTTP 2xx) ?
					codes <- 200
				}
			}

			It("Should not emit messages", func() {
				messages := make(chan *emit.StateChangeMessage)
				codes := make(chan int)
				go emitHttpOk(codes)
				go analyze(resource, codes, messages)
				Consistently(messages).ShouldNot(Receive())
			})
		})

		Context("When getting failures", func() {

			emitCodesFromSlice := func(codeChan chan<- int, codesToEmit []int) {
				for _, code := range codesToEmit {
					codeChan <- code
					time.Sleep(2 * time.Millisecond)
				}
			}

			var messages chan *emit.StateChangeMessage
			var codes chan int
			BeforeEach(func () {
				messages = make(chan *emit.StateChangeMessage)
				codes = make(chan int)
			})

			It("Should emit failure message", func() {
				go emitCodesFromSlice(codes, []int{200, 0, 0})
				go analyze(resource, codes, messages)

				var receivedMessage *emit.StateChangeMessage
				Eventually(messages).Should(Receive(&receivedMessage))
				Expect(receivedMessage.IsOk).To(BeFalse())
			})

			It("Should emit a recovery message when recovering", func() {
				go emitCodesFromSlice(codes, []int{0, 0, 200, 200, 200})
				go analyze(resource, codes, messages)

				var receivedMessage *emit.StateChangeMessage
				Eventually(messages).Should(Receive(&receivedMessage))
				Expect(receivedMessage.IsOk).To(BeFalse())

				Eventually(messages).Should(Receive(&receivedMessage))
				Expect(receivedMessage.IsOk).To(BeTrue())
			})

			It("Should have the correct values in the message", func () {
				go emitCodesFromSlice(codes, []int{200, 0, 0})
				go analyze(resource, codes, messages)

				var receivedMessage *emit.StateChangeMessage
				Eventually(messages).Should(Receive(&receivedMessage))
				Expect(receivedMessage.Resource).To(Equal(resource))
				Expect(receivedMessage.Codes).To(ConsistOf(200, 0, 0))
			})

		})
	})
})
