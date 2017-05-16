package emit

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type fakeEmitter struct {
	hasBeenCalled bool
}

func (emitter *fakeEmitter) Emit(m *StateChangeMessage) {
	emitter.hasBeenCalled = true
}

func (emitter *fakeEmitter) HasBeenCalled() bool {
	return emitter.hasBeenCalled
}

var _ = Describe("Worker", func() {

	Describe("emitMessages", func() {
		It("Should Emit messages on all configured emitters", func() {
			emitters := []emitter{&fakeEmitter{}, &fakeEmitter{}}
			messages := make(chan *StateChangeMessage)

			// emitMessages is blocking, so should be run asynchronously
			go emitMessages(emitters, messages)
			messages <- &StateChangeMessage{}

			Eventually(emitters[0].(*fakeEmitter).HasBeenCalled).Should(BeTrue())
			Eventually(emitters[1].(*fakeEmitter).HasBeenCalled).Should(BeTrue())
		})
	})
})
