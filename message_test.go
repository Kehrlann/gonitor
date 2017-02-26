package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StateChangeMessage", func() {

	Describe("Printing messages", func() {
		var positiveMessage *StateChangeMessage
		var negativeMessage *StateChangeMessage

		BeforeEach(func() {
			positiveMessage = NewMessage("http://test.com", true, []int{1, 2, 3})
			negativeMessage = NewMessage("http://test.com", false, []int{1, 2, 3})
		})

		Context("With any message", func() {
			It("Should print the codes", func() {
				Expect(positiveMessage.String()).To(ContainSubstring("[1 2 3]"))
			})

			It("Should print the url", func() {
				Expect(positiveMessage.String()).To(ContainSubstring("http://test.com"))
			})
		})

		Context("With a positive message", func() {
			It("Should print true", func() {
				Expect(positiveMessage.String()).To(ContainSubstring("true"))
			})
		})

		Context("With a negative message", func() {
			It("Should print false", func() {
				Expect(negativeMessage.String()).To(ContainSubstring("false"))
			})
		})
	})

})
