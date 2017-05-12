package config

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Smtp : ", func() {

	Describe("Smtp.IsValid : ", func() {
		smtp := &Smtp{
			FromAddress: "address@example.com",
			FromName:    "My Name",
			To:          []string{"recipient@example.com"},
			Host:        "host.example.com",
			Username:    "username",
			Password:    "password",
			Port:        25}

		It("Should be valid when everything is filed", func() {
			Expect(smtp.IsValid()).To(BeTrue())
		})

		It("Should be invalid without a from address", func() {
			no_address := *smtp
			no_address.FromAddress = ""
			Expect(no_address.IsValid()).To(BeFalse())
		})

		It("Should be invalid without a to address", func() {
			no_recipient := *smtp
			no_recipient.To = []string{}
			Expect(no_recipient.IsValid()).To(BeFalse())
		})

		It("Should be invalid without a host", func() {
			no_host := *smtp
			no_host.Host = ""
			Expect(no_host.IsValid()).To(BeFalse())
		})

		It("Should be invalid without a port", func() {
			no_port := *smtp
			no_port.Port = 0
			Expect(no_port.IsValid()).To(BeFalse())
		})
	})

	Describe("FormatFromHeader", func() {
		It("Should format the headers correctly", func() {
			smtp := &Smtp{FromAddress: "address@example.com", FromName: "My Name"}
			Expect(smtp.FormatFromHeader()).To(Equal("My Name <address@example.com>"))
		})
	})
})