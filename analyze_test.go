package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Analyze", func() {

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

		resource := Resource{"Url", 2, 2, 10, 3, ""}

		Context("When not polling", func() {
			It("Should not emit messages", func() {
				messages := make(chan *StateChangeMessage)
				codes := make(chan int)
				go Analyze(resource, codes, messages)
				Consistently(messages).ShouldNot(Receive())
			})
		})

		Context("When getting only successes", func() {

			emitHttpOk := func(codes chan<- int) {
				for range time.Tick(2 * time.Millisecond) {
					// TODO : emit random Codes between 200 and 300 (HTTP 2xx)
					codes <- 200
				}
			}

			It("Should not emit messages", func() {
				messages := make(chan *StateChangeMessage)
				codes := make(chan int)
				go emitHttpOk(codes)
				go Analyze(resource, codes, messages)
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

			It("Should emit failure message", func() {
				messages := make(chan *StateChangeMessage)
				codes := make(chan int)
				go emitCodesFromSlice(codes, []int{200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 0, 0, 0})
				go Analyze(resource, codes, messages)

				var receivedMessage *StateChangeMessage
				Eventually(messages).Should(Receive(&receivedMessage))
				Expect(receivedMessage.IsOk).To(BeFalse())
				Expect(receivedMessage.Resource.Url).To(Equal("Url"))
			})

			It("Should emit a recovery message when recovering", func() {
				messages := make(chan *StateChangeMessage)
				codes := make(chan int)
				go emitCodesFromSlice(codes, []int{0, 0, 0, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200})
				go Analyze(resource, codes, messages)

				var receivedMessage *StateChangeMessage
				Eventually(messages).Should(Receive(&receivedMessage))
				Expect(receivedMessage.IsOk).To(BeFalse())
				Expect(receivedMessage.Resource.Url).To(Equal("Url"))

				Eventually(messages).Should(Receive(&receivedMessage))
				Expect(receivedMessage.IsOk).To(BeTrue())

			})

		})
	})
})
