package monitor

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"

	"github.com/jarcoal/httpmock"
)

var _ = Describe("Http", func() {

	log.SetLevel(log.PanicLevel)
	Describe("httpFetcher.fetch", func() {
		Context("When fetching from a website", func() {
			var fetcher *httpFetcher

			BeforeEach(func() {
				fetcher = &httpFetcher{&http.Client{}}
				httpmock.Activate()
			})

			It("Should get the correct code", func() {
				defer httpmock.DeactivateAndReset()

				hasRequestBeenMade := false
				httpmock.RegisterResponder(
					"GET",
					"http://example.com",
					func(req *http.Request) (*http.Response, error) {
						hasRequestBeenMade = true
						return httpmock.NewStringResponse(200, `Woo-hoo`), nil
					})

				status_code := fetcher.fetch("http://example.com")

				Expect(status_code).To(Equal(200))
				Expect(hasRequestBeenMade).To(BeTrue())
			})

			It("Should return 0 when timing out", func () {
				defer httpmock.DeactivateAndReset()
				// Don't register responder -> the request should fail

				status_code := fetcher.fetch("http://example.com")

				Expect(status_code).To(Equal(0))
			})
		})
	})

})
