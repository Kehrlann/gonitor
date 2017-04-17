package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StateChangeMessage : ", func() {

	var positiveMessage *StateChangeMessage
	var negativeMessage *StateChangeMessage

	BeforeEach(func() {
		res := Resource{"http://test.com", 60, 2, 10, 3 }
		positiveMessage = RecoveryMessage(res, []int{1, 2, 3})
		negativeMessage = ErrorMessage(res, []int{1, 2, 3})
	})

	Describe("Printing messages : ", func() {

		Context("With any message", func() {
			It("Should print the Codes", func() {
				Expect(positiveMessage.String()).To(ContainSubstring("[1 2 3]"))
			})

			It("Should print the Url", func() {
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

	Describe("E-mail helpers : ", func() {

		Describe("Subject", func() {

			Context("With a positive message", func() {
				It("Should contains \"recovery\"", func() {
					Expect(positiveMessage.MailSubject()).To(Equal("[Gonitor] Recovery for 'http://test.com'"))
				})
			})

			Context("With a negative message", func() {
				It("Should contain \"failure\"", func() {
					Expect(negativeMessage.MailSubject()).To(Equal("[Gonitor] Failure for 'http://test.com'"))
				})
			})

		})

		Describe("Body", func() {

			Context("With a positive message", func() {
				It("Should contains \"recovery\"", func() {
					Expect(positiveMessage.MailBody()).To(
						Equal(
							"Hi !<br><br>" +
								"" +
								"This is an automated message from Gonitor.<br>" +
								"It seems <strong style=\"color:#5cb85c\">http://test.com</strong> has recovered.<br>" +
								"The following HTTP codes were received : [1 2 3].<br><br>" +
								"The config used is :" +
								"<ul>" +
								"<li>Interval : 60 seconds</li>" +
								"<li>Number of tries : 10</li>" +
								"<li>Failure threshold : 3</li>" +
								"<li>Recovery threshold : 10</li>" +
								"<li>Timeout : 2 seconds</li>" +
								"</ul>"))
				})
			})

			Context("With a negative message", func() {
				It("Should contain \"failure\"", func() {
					Expect(negativeMessage.MailBody()).To(
						Equal(
							"Hi !<br><br>" +
								"" +
								"This is an automated message from Gonitor.<br>" +
								"It seems an error occurred when polling <strong style=\"color:#d9534f\">http://test.com</strong>.<br>" +
								"The following HTTP codes were received : [1 2 3].<br><br>" +
								"The config used is :" +
								"<ul>" +
								"<li>Interval : 60 seconds</li>" +
								"<li>Number of tries : 10</li>" +
								"<li>Failure threshold : 3</li>" +
								"<li>Recovery threshold : 10</li>" +
								"<li>Timeout : 2 seconds</li>" +
								"</ul>"))
				})
			})

		})
	})

})
