package actor

import (
	"reflect"
	"sync"

	"github.com/ZHANG-JIHUI/zephyr/tools/log"
)

type DeadLetter struct {
	eventStream *EventStream
	pid         *PID
}

func NewDeadLetter(eventStream *EventStream) *DeadLetter {
	return &DeadLetter{
		eventStream: eventStream,
		pid:         NewPID(LocalLookupAddr, "deadLetter"),
	}
}

func (slf *DeadLetter) Start() {
	//TODO implement me
	panic("implement me")
}

func (slf *DeadLetter) PID() *PID {
	return slf.pid
}

func (slf *DeadLetter) Send(target *PID, message any, sender *PID) {
	log.Warn("[ACTOR] dead letter send", log.Any("target", target),
		log.Any("message", reflect.TypeOf(message)), log.Any("sender", sender))
	slf.eventStream.Publish(&EventDeadLetter{
		Target:  target,
		Message: message,
		Sender:  sender,
	})
}

func (slf *DeadLetter) Invoke(envelopes []Envelope) {
	//TODO implement me
	panic("implement me")
}

func (slf *DeadLetter) Shutdown(wg *sync.WaitGroup) {
	//TODO implement me
	panic("implement me")
}
