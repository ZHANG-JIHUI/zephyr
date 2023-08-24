package actor

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/ZHANG-JIHUI/zephyr/tools/log"
)

type Processor interface {
	Start()
	PID() *PID
	Send(receiver *PID, message any, sender *PID)
	Invoke(envelopes []Envelope)
	Shutdown(wg *sync.WaitGroup)
}

type processor struct {
	Options
	inbox    *Inbox
	ctx      *Context
	pid      *PID
	restarts int32
	buffer   []Envelope
}

func newProcessor(engine *Engine, opts Options) *processor {
	pid := NewPID(engine.address, opts.Name, opts.Tags...)
	ctx := newContext(engine, pid)
	proc := &processor{
		Options: opts,
		inbox:   NewInbox(opts.InboxSize),
		ctx:     ctx,
		pid:     pid,
	}
	proc.inbox.Start(proc)
	return proc
}

func (slf *processor) Start() {
	receiver := slf.Producer()
	slf.ctx.receiver = receiver
	slf.ctx.message = Started{}
	applyMiddleware(receiver.Receive, slf.Options.Middleware...)(slf.ctx)
	slf.ctx.engine.EventStream.Publish(&EventActivation{PID: slf.pid})
	log.Debug("[ACTOR] processor started", log.Any("PID", slf.pid))
	if len(slf.buffer) > 0 {
		slf.Invoke(slf.buffer)
		slf.buffer = nil
	}
}

func (slf *processor) PID() *PID { return slf.pid }

func (slf *processor) Send(_ *PID, message any, sender *PID) {
	slf.inbox.Send(Envelope{Message: message, Sender: sender})
}

func (slf *processor) Invoke(envelopes []Envelope) {

	var (
		nEnvelope  = len(envelopes)
		nProcessed = 0
	)
	defer func() {
		if v := recover(); v != nil {
			slf.ctx.message = Stopped{}
			slf.ctx.receiver.Receive(slf.ctx)
			slf.buffer = make([]Envelope, nEnvelope-nProcessed)
			for i := 0; i < nEnvelope-nProcessed; i++ {
				slf.buffer[i] = envelopes[i+nProcessed]
			}
			if slf.Options.MaxRestarts > 0 {
				slf.tryRestart(v)
			}
		}
	}()
	for i := 0; i < nEnvelope; i++ {
		nProcessed++
		envelope := envelopes[i]
		if pill, ok := envelope.Message.(poisonPill); ok {
			slf.cleanup(pill.wg)
			return
		}
		slf.ctx.message = envelope.Message
		slf.ctx.sender = envelope.Sender
		receiver := slf.ctx.receiver
		if len(slf.Options.Middleware) > 0 {
			applyMiddleware(receiver.Receive, slf.Options.Middleware...)(slf.ctx)
		} else {
			receiver.Receive(slf.ctx)
		}
	}
}

func (slf *processor) Shutdown(wg *sync.WaitGroup) { slf.cleanup(wg) }

func (slf *processor) cleanup(wg *sync.WaitGroup) {
	_ = slf.inbox.Stop()
	slf.ctx.engine.Registry.Remove(slf.pid)
	slf.ctx.message = Stopped{}
	applyMiddleware(slf.ctx.receiver.Receive, slf.Options.Middleware...)(slf.ctx)
	if slf.ctx.parent != nil {
		slf.ctx.parent.children.Delete(slf.Name)
	}
	if slf.ctx.children.Len() > 0 {
		children := slf.ctx.Children()
		for _, pid := range children {
			if wg != nil {
				wg.Add(1)
			}
			proc := slf.ctx.engine.Registry.get(pid)
			proc.Shutdown(wg)
		}
	}
	log.Debug("[ACTOR] processor cleanup", log.Any("pid", slf.pid))
	slf.ctx.engine.EventStream.Publish(&EventTermination{PID: slf.pid})
	if wg != nil {
		wg.Done()
	}
}

func (slf *processor) tryRestart(v any) {
	slf.restarts++
	if msg, ok := v.(*InternalError); ok {
		log.Error("[ACTOR] processor restart failed", log.String("from", msg.From), log.Err(msg.Err))
		time.Sleep(slf.Options.RestartDelay)
		slf.Start()
		return
	}
	fmt.Println(string(debug.Stack()))
	if slf.restarts == slf.MaxRestarts {
		log.Error("[ACTOR] processor max restarts exceeded, shutting down...",
			log.Any("pid", slf.pid), log.Int32("restarts", slf.restarts))
		slf.cleanup(nil)
		return
	}
	log.Error("[ACTOR] processor restarting", log.Int32("restarts", slf.restarts),
		log.Int32("maxRestarts", slf.MaxRestarts), log.Any("pid", slf.pid), log.Any("reason", v))
	time.Sleep(slf.Options.RestartDelay)
	slf.Start()
}

func applyMiddleware(rcv ReceiveFunc, middleware ...MiddlewareFunc) ReceiveFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		rcv = middleware[i](rcv)
	}
	return rcv
}
