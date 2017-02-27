package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Analyze", func() {

	Describe("Computing state", func() {

		Context("Given all successes", func() {
			codesToAnalyze := []int{200, 200, 200, 200, 200, 200, 200, 200, 200, 200}

			It("Should not be a failure", func() {
				failure, _ := computeState(codesToAnalyze)
				Expect(failure).To(BeFalse())
			})

			It("Should allow recovery", func() {
				_, canRecover := computeState(codesToAnalyze)
				Expect(canRecover).To(BeTrue())
			})
		})

		Context("Given one failure", func() {

			codesToAnalyze := []int{200, 200, 200, 200, 200, 200, 200, 200, 200, 0}
			It("Should not be a failure", func() {
				failure, _ := computeState(codesToAnalyze)
				Expect(failure).To(BeFalse())
			})

			It("Should not allow recovery", func() {
				_, canRecover := computeState(codesToAnalyze)
				Expect(canRecover).To(BeFalse())
			})
		})

		Context("Given all failures", func() {
			codesToAnalyze := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

			It("Should be a failure", func() {
				failure, _ := computeState(codesToAnalyze)
				Expect(failure).To(BeTrue())
			})

			It("Should not allow recovery", func() {
				_, canRecover := computeState(codesToAnalyze)
				Expect(canRecover).To(BeFalse())
			})
		})

	})


})
