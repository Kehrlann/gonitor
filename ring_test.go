package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"container/ring"
)

var _ = Describe("Ring", func() {

	Describe("Ring to Slice", func() {
		Context("With a normal ring", func() {
			var r *ring.Ring
			var s []int

			BeforeEach(func() {
				r = ring.New(5)
				s = make([]int, 5)
				for i := 0; i < 5; i++{
					r.Value = i +1
					r = r.Next()
					s[i] = i + 1
				}
			})

			It("Should dump the values in order", func() {
				intSlice := RingToIntSlice(r)
				Expect(intSlice).To(Equal(s))
			})

			It("Should not mutate the ring", func() {
				RingToIntSlice(r)
				Expect(r.Value).To(Equal(1))
				r = r.Next()
				Expect(r.Value).To(Equal(2))
			})
		})

		Context("With a half-full ring", func() {
			var r *ring.Ring
			var s []int

			BeforeEach(func() {
				r = ring.New(10)
				s = make([]int, 10)
				for i := 0; i < 5; i++{
					r.Value = i +1
					r = r.Next()
					s[i] = i + 1
				}
			})

			It("Should only keep actual ring values", func() {
				intSlice := RingToIntSlice(r)
				Expect(intSlice).To(Equal(s))
			})
		})



		Context("With an empty ring", func() {
			var r *ring.Ring
			var s []int

			BeforeEach(func() {
				r = ring.New(10)
				s = make([]int, 10)
			})

			It("Should return a zero-ed slice", func() {
				intSlice := RingToIntSlice(r)
				Expect(intSlice).To(Equal(s))
			})
		})
	})
})
