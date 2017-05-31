package alert

import (
	"github.com/kehrlann/gonitor/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"math"
)

var _ = Describe("websockets -> ", func() {
	var message *StateChangeMessage
	var mailer *FakeMailer

	BeforeEach(func() {
		res := config.Resource{"http://test.com", 60, 2, 10, 3, "" }
		// TODO : use pointers to resources ?
		message = RecoveryMessage(res, []int{1, 2, 3})
		mailer = &FakeMailer{}
	})

	Context("Tech tests ", func() {
		It("is fun", func() {
			var a uint8
			a = math.MaxUint8

			Expect(a).To(BeNumerically("==", 255))

			b := a + 1
			Expect(b).To(BeNumerically("==", 0))
		})
	})
})
