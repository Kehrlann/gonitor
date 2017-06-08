package alert

import (
	"github.com/kehrlann/gonitor/config"

	. "github.com/onsi/ginkgo"
	"github.com/kehrlann/gonitor/monitor"
)

var _ = Describe("websockets -> ", func() {
	var message *monitor.StateChangeMessage
	// TODO : test me !
	// TODO : you need to read the message to be able to know whether the connection was closed or not ... test that
	// TODO : test timeout on client
	// TODO : test connection closed
	//			Caveat : apparently you have to read messages to be sure that the connection has been closed
	// 				conn.SetReadDeadline(time.Now().Add(time.Second))
	//				conn.ReadMessage()

	BeforeEach(func() {
		res := config.Resource{"http://test.com", 60, 2, 10, 3, "" }
		// TODO : use pointers to resources ?
		message = monitor.RecoveryMessage(res, []int{1, 2, 3})
	})

	Context("Tech tests ", func() {
		It("is fun", func() {
		})
	})
})
