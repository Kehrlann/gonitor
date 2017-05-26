package web

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
)

var _ = Describe("Server", func() {

	var cleanup func()

	BeforeSuite(func() {
		cleanup = Serve()
	})

	AfterSuite(func() {
		if cleanup != nil {
			cleanup()
		}
	})

	It("should serve an index page with a code 200", func() {
		client := &http.Client{}

		getIndex := func() int {
			r, err := client.Get("http://127.0.0.1:3000/")
			if err != nil {
				return 0
			}

			return r.StatusCode
		}

		Eventually(getIndex).Should(Equal(200))
	})

	It("should serve an index with the word `Websocket` in it", func() {
		client := &http.Client{}

		getBody := func() (string, error) {
			r, err := client.Get("http://127.0.0.1:3000/")
			if err != nil {
				return "", err
			}
			defer r.Body.Close()

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return "", err
			}

			return string(body), nil
		}

		Eventually(getBody).Should(ContainSubstring("Websocket"))
	})
})
