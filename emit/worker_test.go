package emit

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// since the Emit method is called asynchronously, we use a chan of messages to denoted that it has been called ;
// instead of a simple boolean (which causes a data-race)
type fakeEmitter struct {
	calls chan *StateChangeMessage
}

func (emitter *fakeEmitter) Emit(m *StateChangeMessage) {
	emitter.calls <- m
}

func newFake() *fakeEmitter{
	return &fakeEmitter{make(chan *StateChangeMessage)}
}

var _ = Describe("Worker", func() {

	Describe("emitMessages", func() {
		It("Should Emit messages on all configured emitters", func() {
			emitters := []emitter{newFake(), newFake()}
			messages := make(chan *StateChangeMessage, 1)
			messages <- &StateChangeMessage{}
			close(messages)

			// emitMessages is blocking, so should be run asynchronously
			emitMessages(emitters, messages)

			Eventually(emitters[0].(*fakeEmitter).calls).Should(Receive())
			Eventually(emitters[1].(*fakeEmitter).calls).Should(Receive())
		})
	})
})
