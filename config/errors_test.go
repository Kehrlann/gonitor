package config

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration : ", func() {
	Describe("NoDefaultConfigError : ", func() {
		It("Should be an error type with the correct message", func() {
			err := NewDefaultConfigError()

			Expect(err.Error()).To(Equal("No default config found at ./gonitor.config.json"))
			Expect(err.HelpMessage).To(ContainSubstring(DEFAULT_CONFIG_PATH))
		})
	})
})