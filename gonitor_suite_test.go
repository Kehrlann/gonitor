package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGonitor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gonitor Suite")
}
