package config

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGonitor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gonitor/config Suite")
}
