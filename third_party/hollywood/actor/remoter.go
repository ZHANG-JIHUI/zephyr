package actor

type Remoter interface {
	Start()
	Send(target *PID, msg any, sender *PID)
	Address() string
}
