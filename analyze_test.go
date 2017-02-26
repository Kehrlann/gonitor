package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


var _ = Describe("Analyze", func() {

	Describe("Computing state", func() {

		Context("Given all successes", func() {

			It("Should not be a failure", func() {
				codesToAnalyze := []int{200, 200, 200, 200, 200, 200, 200, 200, 200, 200}
				failure, _ := computeState(codesToAnalyze)
				Expect(failure).To(BeFalse())
			})

			It("Should allow recovery", func() {
				codesToAnalyze := []int{200, 200, 200, 200, 200, 200, 200, 200, 200, 200}
				_, canRecover := computeState(codesToAnalyze)
				Expect(canRecover).To(BeTrue())
			})
		})

		Context("Given one failure", func() {

			It("Should not be a failure", func() {
				codesToAnalyze := []int{200, 200, 200, 200, 200, 200, 200, 200, 200, 0}
				failure, _ := computeState(codesToAnalyze)
				Expect(failure).To(BeFalse())
			})

			It("Should not allow recovery", func() {
				codesToAnalyze := []int{200, 200, 200, 200, 200, 200, 200, 200, 200, 0}
				_, canRecover := computeState(codesToAnalyze)
				Expect(canRecover).To(BeFalse())
			})
		})


		Context("Given all failures", func() {

			It("Should be a failure", func() {
				codesToAnalyze := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
				failure, _ := computeState(codesToAnalyze)
				Expect(failure).To(BeTrue())
			})

			It("Should not allow recovery", func() {
				codesToAnalyze := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
				_, canRecover := computeState(codesToAnalyze)
				Expect(canRecover).To(BeFalse())
			})
		})

	})

})
