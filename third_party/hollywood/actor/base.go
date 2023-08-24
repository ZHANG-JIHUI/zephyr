package actor

type Base interface {
	PID() *PID
	Producer() Producer
}
