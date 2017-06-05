package web

import (
	"net/http"
	"io/ioutil"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/gorilla/websocket"
	log "github.com/Sirupsen/logrus"
	testlog "github.com/Sirupsen/logrus/hooks/test"
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

		// TODO : cleanup (otherwise the server is still open and this connection remains open
		It("Should serve the websocket upgrade page", func() {
			getResp := func() (string) {
				dialer := &websocket.Dialer{}
				conn, _, _ := dialer.Dial("ws://127.0.0.1:3000/ws", nil)
				_, message, _ := conn.ReadMessage()
				return string(message)
			}

			Eventually(getResp).Should(ContainSubstring("date"))
		})

		It("Should close the connection", func () {
			dialer := &websocket.Dialer{}
			conn, _, _ := dialer.Dial("ws://127.0.0.1:3000/ws", nil)
			time.Sleep(1 * time.Second)
			conn.WriteControl(websocket.CloseMessage, nil, time.Now().Add(time.Second))
			//err := conn.Close()
			time.Sleep(5 * time.Second)

			//Expect(err).To(BeNil())
		})
	})

})
