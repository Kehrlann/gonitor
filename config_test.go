package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

var _ = Describe("Config : ", func() {

	Describe("LoadConfig : ", func() {
		Context("When loading a valid config from a file", func() {
			It("Should return the proper resources", func() {
				// Arrange
				tempFile := createConfigFile(`{
				"resources" : 	[
									{
										"url" : 				"http://www.example.com",
										"intervalInSeconds" : 	60,
										"timeoutInSeconds" :	2,
										"numberOfTries" :		10,
										"failureThreshold" :	3
									},
									{
										"url" : 				"http://www.example.test",
										"intervalInSeconds" : 	120,
										"timeoutInSeconds" :	10,
										"numberOfTries" :		10,
										"failureThreshold" :	10
									}
								]
			}`)
				defer os.Remove(tempFile.Name())

				// Act
				config, err := LoadConfig(tempFile.Name())

				// Assert
				Expect(err).To(BeNil())
				Expect(config).ToNot(BeNil())
				Expect(len(config.Resources)).To(Equal(2))

				first := config.Resources[0]
				resource := Resource{"http://www.example.com", 60, 2, 10, 3}
				Expect(first).To(Equal(resource))

				second := config.Resources[1]
				resource = Resource{"http://www.example.test", 120, 10, 10, 10}
				Expect(second).To(Equal(resource))
			})

			It("Should load the SMTP config when it exists", func() {
				tempFile := createConfigFile(`{
				"smtp" :	{
								"host" 			: "smtp.example.com",
								"port"			: 25,
								"username" 		: "user",
								"password" 		: "password123",
								"fromaddress" 	: "address@example.com",
								"fromname" 		: "Mr Example",
								"to" 			: ["recipient@example.com", "admin@example.com"]
							}
			}`)
				defer os.Remove(tempFile.Name())

				// Act
				config, err := LoadConfig(tempFile.Name())

				// Assert
				Expect(err).To(BeNil())
				Expect(config).ToNot(BeNil())

				smtp := config.Smtp
				Expect(smtp.Host).To(Equal("smtp.example.com"))
				Expect(smtp.Port).To(Equal(25))
				Expect(smtp.FromAddress).To(Equal("address@example.com"))
				Expect(smtp.FromName).To(Equal("Mr Example"))
				Expect(smtp.Username).To(Equal("user"))
				Expect(smtp.Password).To(Equal("password123"))
				Expect(smtp.To).To(ConsistOf("recipient@example.com", "admin@example.com"))

			})
		})

		Context("When loading an invalid config from a file", func() {
			It("Should throw an error", func() {
				tempfile, _ := ioutil.TempFile("", "gonitor.config.json")
				defer os.Remove(tempfile.Name())
				tempfile.WriteString("hello i am invalid !")
				tempfile.Close()

				config, err := LoadConfig(tempfile.Name())

				Expect(err).ToNot(BeNil())
				Expect(config).To(BeNil())
			})
		})

		Context("When trying to load a non-existing file", func() {
			It("Should throw an error", func() {
				config, err := LoadConfig("i_am_not_a_file.json")

				Expect(err).ToNot(BeNil())
				Expect(config).To(BeNil())
			})
		})
	})

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

		It("Should be invalid without a from address", func(){
			no_address  := *smtp
			no_address.FromAddress = ""
			Expect(no_address.IsValid()).To(BeFalse())
		})

		It("Should be invalid without a to address", func(){
			no_recipient := *smtp
			no_recipient.To = []string{}
			Expect(no_recipient.IsValid()).To(BeFalse())
		})

		It("Should be invalid without a host", func(){
			no_host := *smtp
			no_host.Host = ""
			Expect(no_host.IsValid()).To(BeFalse())
		})

		It("Should be invalid without a port", func(){
			no_port := *smtp
			no_port.Port = 0
			Expect(no_port.IsValid()).To(BeFalse())
		})
	})
})

func createConfigFile(config string) *os.File {
	tempfile, _ := ioutil.TempFile("", "gonitor.config.json")
	tempfile.WriteString(config)
	tempfile.Close()
	return tempfile
}
