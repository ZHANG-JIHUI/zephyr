package actor

type Envelope struct {
	Message any
	Sender  *PID
}

func (slf *Envelope) GetMessage() any {
	return slf.Message
}

func (slf *Envelope) GetSender() *PID {
	return slf.Sender
}
