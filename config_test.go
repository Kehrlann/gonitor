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
			tempfile.Close()

			Expect(tempfile.Name()).ToNot(BeNil())
		})
	})

	Context("When loading an invalid config from a file", func() {
		It("Should throw an error", func() {

		})
	})

	Context("When trying to load a non-existing file", func() {
		It("Should throw an error", func() {

		})
	})
})
