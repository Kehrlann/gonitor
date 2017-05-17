package alert

type emitter interface {
	Emit(message *StateChangeMessage)
}
