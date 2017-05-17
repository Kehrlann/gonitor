package emit

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sync"
)

var _ = Describe("Worker", func() {

	Describe("emitMessages", func() {
		It("Should Emit messages on all configured emitters", func() {
			emitters := []emitter{newFakeChan(), newFakeMutex()}
			messages := make(chan *StateChangeMessage, 1)
			messages <- &StateChangeMessage{}
			close(messages)

			// emitMessages is blocking, so should be run asynchronously
			emitMessages(emitters, messages)

			Eventually(emitters[0].(*fakeChanEmitter).calls).Should(Receive())
			Eventually(emitters[1].(*fakeMutexEmitter).HasBeenCalled).Should(BeTrue())
		})
	})
})

type fakeChanEmitter struct {
	calls chan *StateChangeMessage
}

func (emitter *fakeChanEmitter) Emit(m *StateChangeMessage) {
	emitter.calls <- m
}

func newFakeChan() *fakeChanEmitter {
	return &fakeChanEmitter{make(chan *StateChangeMessage)}
}

type fakeMutexEmitter struct {
	hasBeenCalled bool
	lock sync.Mutex
}

func (emitter *fakeMutexEmitter) Emit(m *StateChangeMessage) {
	emitter.lock.Lock()
	defer emitter.lock.Unlock()
	emitter.hasBeenCalled = true
}

func (emitter *fakeMutexEmitter) HasBeenCalled() bool {
	emitter.lock.Lock()
	defer emitter.lock.Unlock()
	return emitter.hasBeenCalled
}

func newFakeMutex() *fakeMutexEmitter {
	return &fakeMutexEmitter{false, sync.Mutex{}}
}
