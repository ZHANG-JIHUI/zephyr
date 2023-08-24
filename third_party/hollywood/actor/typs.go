package actor

import "sync"

type Started struct{}
type Stopped struct{}

type (
	EventDeadLetter struct {
		Target  *PID
		Message any
		Sender  *PID
	}
	EventActivation  struct{ PID *PID }
	EventTermination struct{ PID *PID }
)

type InternalError struct {
	From string
	Err  error
}

type poisonPill struct {
	wg *sync.WaitGroup
}
