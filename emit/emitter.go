package emit

type emitter interface {
	Emit(message *StateChangeMessage)
}
