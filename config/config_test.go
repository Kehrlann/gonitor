package config

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("Configuration : ", func() {
	log.SetLevel(log.PanicLevel)
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
										"failureThreshold" :	3,
										"command" :				"foobar"
									},
									{
										"url" : 				"http://www.example.test",
										"intervalInSeconds" : 	120,
										"timeoutInSeconds" :	10,
										"numberOfTries" :		10,
										"failureThreshold" :	10,
										"command" :				null
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
				resource := Resource{"http://www.example.com", 60, 2, 10, 3, "foobar"}
				Expect(first).To(Equal(resource))

				second := config.Resources[1]
				resource = Resource{"http://www.example.test", 120, 10, 10, 10, ""}
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

			It("Should return the global command when set", func() {
				tempFile := createConfigFile(`{"globalCommand" : "hello.sh"}`)
				defer os.Remove(tempFile.Name())

				// Act
				config, err := LoadConfig(tempFile.Name())

				// Assert
				Expect(err).To(BeNil())
				Expect(config.GlobalCommand).To(Equal("hello.sh"))
			})

			It("Should have an empty string as global command when null or not set", func() {
				// Arrange
				notSetFile := createConfigFile(`{}`)
				nullFile := createConfigFile(`{ "globalCommand" : null }`)
				defer os.Remove(nullFile.Name())
				defer os.Remove(notSetFile.Name())

				// Act
				configNotSet, err := LoadConfig(notSetFile.Name())
				configNull, err2 := LoadConfig(nullFile.Name())

				// Assert
				Expect(err).To(BeNil())
				Expect(configNotSet.GlobalCommand).To(BeEmpty())

				Expect(err2).To(BeNil())
				Expect(configNull.GlobalCommand).To(BeEmpty())
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

		Context("When trying to load a default config that doesn't exist", func() {
			It("Should throw a NoDefaultConfigError", func() {
				temp := DEFAULT_CONFIG_PATH
				DEFAULT_CONFIG_PATH := "i_am_not_a_file.json"
				defer func() { DEFAULT_CONFIG_PATH = temp }()

				_, err := LoadConfig(DEFAULT_CONFIG_PATH)

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).ToNot(BeNil())
			})
		})
	})
})

func createConfigFile(config string) *os.File {
	tempfile, _ := ioutil.TempFile("", "gonitor.config.json")
	tempfile.WriteString(config)
	tempfile.Close()
	return tempfile
}
