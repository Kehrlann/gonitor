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
			tempfile.Close()

			config, err := LoadConfig(tempfile.Name())

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
	})

	Context("When loading an invalid config from a file", func() {
		It("Should throw an error", func() {
			tempfile, _ := ioutil.TempFile("", "config.json")
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
