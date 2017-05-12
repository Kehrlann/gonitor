package emit

type Emitter interface {
	Emit(message *StateChangeMessage)
}
