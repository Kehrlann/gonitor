package main

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gopkg.in/jarcoal/httpmock.v1"
)

var _ = Describe("Http", func() {

	Describe("Fetch", func() {

		Context("When fetching from a website", func() {

			It("Should get the correct code", func() {

				hasRequestBeenMade := false
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()
				httpmock.RegisterResponder(
					"GET",
					"http://example.com",
					func(req *http.Request) (*http.Response, error) {
						hasRequestBeenMade = true
						return httpmock.NewStringResponse(200, `Woo-hoo`), nil
					})

				status_code := fetch(&http.Client{}, "http://example.com")

				Expect(status_code).To(Equal(200))
				Expect(hasRequestBeenMade).To(BeTrue())
			})
		})
	})

})
