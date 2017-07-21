package web

import (
	"net/http"
	"io/ioutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	testlog "github.com/sirupsen/logrus/hooks/test"
	"github.com/kehrlann/gonitor/server/web/handlers"
	"github.com/kehrlann/gonitor/websockets"
	"github.com/gorilla/websocket"
)

var _ = Describe("Server", func() {

	var cleanup func()

	BeforeSuite(func() {
		connections := make(chan websockets.Connection, 10)
		cleanup = serve(handlers.WebsocketHandler{connections})
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
			log.SetOutput(&nilWriter{})
			getCode := func() (int) {
				client := &http.Client{}
				resp, _ := client.Get("http://127.0.0.1:3000/ws")

				return resp.StatusCode
			}

			// act / assert
			Consistently(getCode).Should(Equal(400))
			Eventually(func() []*log.Entry { return hook.Entries }).ShouldNot(BeEmpty())
			Expect(hook.LastEntry().Level).To(Equal(log.ErrorLevel))
		})

		It("Should answer 101 to a ws-enabled client", func () {
			getCode := func() (int, error) {
				dialer := &websocket.Dialer{}
				_, resp, err := dialer.Dial("ws://127.0.0.1:3000/ws", nil)
				return resp.StatusCode, err
			}

			Consistently(getCode).Should(Equal(101))
		})
	})
})


type nilWriter struct {}

func (w *nilWriter) Write(p []byte) (int, error) {
	return 0, nil
}

