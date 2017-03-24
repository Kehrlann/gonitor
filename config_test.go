package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

var _ = Describe("Config", func() {

	Context("When loading a valid config from a file", func() {
		It("Should return the proper data structure", func() {
			tempfile, _ := ioutil.TempFile("", "config.json")
			defer os.Remove(tempfile.Name())
			tempfile.WriteString(
				`{
						"resources" : 	[
											{
												"url" : 				"http://www.example.com",
												"intervalInSeconds" : 	60,
												"timeoutInSeconds" :	2,
												"numberOfTries" :		10,
												"failuteThreshold" :	3
											},
											{
												"url" : 				"http://www.example.test",
												"intervalInSeconds" : 	120,
												"timeoutInSeconds" :	10,
												"numberOfTries" :		10,
												"failuteThreshold" :	10
											}
										]
					}`)
			tempfile.Close()

			config, error := LoadConfig(tempfile.Name())

			Expect(error).To(BeNil())
			Expect(config).ToNot(BeNil())
		})
	})

	Context("When loading an invalid config from a file", func() {
		It("Should throw an error", func() {
			tempfile, _ := ioutil.TempFile("", "config.json")
			defer os.Remove(tempfile.Name())
			tempfile.WriteString("hello i am invalid !")
			tempfile.Close()

			config, error := LoadConfig(tempfile.Name())

			Expect(error).ToNot(BeNil())
			Expect(config).To(BeNil())
		})
	})

	Context("When trying to load a non-existing file", func() {
		It("Should throw an error", func() {
			config, error := LoadConfig("i_am_not_a_file.json")

			Expect(error).ToNot(BeNil())
			Expect(config).To(BeNil())
		})
	})
})
