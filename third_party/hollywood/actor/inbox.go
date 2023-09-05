package actor

import (
	"runtime"
	"sync/atomic"

	"github.com/ZHANG-JIHUI/zephyr/tools/ring"
	"golang.org/x/exp/constraints"
)

const (
	IDLE int32 = iota
	RUNNING
)

type Inbox struct {
	buffer    *ring.Ring[Envelope]
	processor Processor
	scheduler Scheduler
	status    int32
}

func NewInbox[T constraints.Signed](size T) *Inbox {
	return &Inbox{
		buffer:    ring.New[Envelope](int64(size)),
		scheduler: NewScheduler(defaultThroughput),
	}
}

func (slf *Inbox) Start(processor Processor) {
	slf.processor = processor
}

func (slf *Inbox) Stop() error {
	return nil
}

func (slf *Inbox) Send(env Envelope) {
	slf.buffer.Push(env)
	slf.schedule()
}

func (slf *Inbox) schedule() {
	if atomic.CompareAndSwapInt32(&slf.status, IDLE, RUNNING) {
		slf.scheduler.Schedule(slf.process)
	}
}

func (slf *Inbox) process() {
	slf.run()
	atomic.StoreInt32(&slf.status, IDLE)
}

func (slf *Inbox) run() {
	count, throughput := 0, slf.scheduler.Throughput()
	for {
		if count > throughput {
			count = 0
			runtime.Gosched()
		}
		count++
		if envelope, ok := slf.buffer.Pop(); ok {
			slf.processor.Invoke([]Envelope{envelope})
		} else {
			return
		}
	}
}
