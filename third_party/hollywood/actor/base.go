package actor

type Base interface {
	PID() *PID
	Context() *Context
	Producer() Producer
}
