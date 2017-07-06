package web

import (
	"net/http"
	"io/ioutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/Sirupsen/logrus"
	testlog "github.com/Sirupsen/logrus/hooks/test"
	"github.com/kehrlann/gonitor/server/web/handlers"
	"github.com/kehrlann/gonitor/websockets"
)

var _ = Describe("Server", func() {

	var cleanup func()

	BeforeSuite(func() {
		cleanup = serve(handlers.WebsocketHandler{make(chan websockets.Connection, 10)})
	})

	AfterSuite(func() {
		if cleanup != nil {
			cleanup()
		}
	})

	Describe("Index -> ", func() {
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

	Describe("Websocket connexion ->", func() {

		It("Should answer 400 to a non-ws-enabled client, and log an error", func() {
			// arrange
			hook := testlog.NewGlobal()
			log.SetLevel(log.ErrorLevel)
			getCode := func() (int) {
				client := &http.Client{}
				resp, _ := client.Get("http://127.0.0.1:3000/ws")

				return resp.StatusCode
			}

			// act / assert
			Consistently(getCode).Should(Equal(400))
			Expect(hook.Entries).ToNot(BeEmpty())
			Expect(hook.LastEntry().Level).To(Equal(log.ErrorLevel))
		})

		// TODO : add tests for websockets connections being added
	})
})
