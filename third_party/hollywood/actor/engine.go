package actor

import (
	"sync"
	"time"

	"github.com/ZHANG-JIHUI/zephyr/tools/log"
)

type (
	Engine struct {
		address     string
		remoter     Remoter
		deadLetter  Processor
		EventStream *EventStream
		Registry    *Registry
	}

	Config struct {
		PIDSeparator string
	}

	Receiver interface {
		Receive(*Context)
	}

	Producer func() Receiver

	funcReceiver struct {
		fn func(ctx *Context)
	}
)

func NewEngine(config ...Config) *Engine {
	engine := &Engine{
		EventStream: NewEventStream(),
		address:     LocalLookupAddr,
	}
	if len(config) == 1 {
		engine.configure(config[0])
	}
	engine.Registry = newRegistry(engine)
	engine.deadLetter = NewDeadLetter(engine.EventStream)
	engine.Registry.add(engine.deadLetter)
	return engine
}

func (slf *Engine) configure(config Config) {
	if config.PIDSeparator != "" {
		pidSeparator = config.PIDSeparator
	}
}

func (slf *Engine) WithRemoter(remoter Remoter) {
	slf.remoter = remoter
	slf.address = remoter.Address()
	remoter.Start()
}

func (slf *Engine) Spawn(producer Producer, name string, opts ...Option) *PID {
	options := DefaultOpts(producer)
	options.Name = name
	for _, opt := range opts {
		opt(&options)
	}
	proc := newProcessor(slf, options)
	return slf.SpawnProc(proc)
}

func (slf *Engine) SpawnFunc(fn func(*Context), id string, opts ...Option) *PID {
	return slf.Spawn(newFuncReceiver(fn), id, opts...)
}

func (slf *Engine) SpawnProc(proc Processor) *PID {
	slf.Registry.add(proc)
	proc.Start()
	return proc.PID()
}

func (slf *Engine) Address() string {
	return slf.address
}

func (slf *Engine) Request(pid *PID, message any, timeout time.Duration) *Response {
	resp := NewResponse(slf, timeout)
	slf.Registry.add(resp)
	slf.SendWithSender(pid, message, resp.PID())
	return resp
}

func (slf *Engine) SendWithSender(pid *PID, msg any, sender *PID) {
	slf.send(pid, msg, sender)
}

func (slf *Engine) Send(pid *PID, msg any) {
	slf.send(pid, msg, nil)
}

func (slf *Engine) SendRepeat(pid *PID, msg any, interval time.Duration) SendRepeater {
	clonedPID := *pid.CloneVT()
	repeater := SendRepeater{
		engine:   slf,
		self:     nil,
		target:   &clonedPID,
		msg:      msg,
		interval: interval,
		cancel:   make(chan struct{}, 1),
	}
	repeater.start()
	return repeater
}

func (slf *Engine) send(pid *PID, msg any, sender *PID) {
	if slf.isLocalMessage(pid) {
		slf.SendLocal(pid, msg, sender)
		return
	}
	if slf.remoter == nil {
		log.Error("[ACTOR] failed to send message because remoter is nil")
		return
	}
	slf.remoter.Send(pid, msg, sender)
}

func (slf *Engine) Poison(pid *PID, wg ...*sync.WaitGroup) {
	var _wg *sync.WaitGroup
	if len(wg) > 0 {
		_wg = wg[0]
		_wg.Add(1)
	}
	proc := slf.Registry.get(pid)
	if proc != nil {
		slf.SendLocal(pid, poisonPill{wg: _wg}, nil)
	}
}

func (slf *Engine) SendLocal(target *PID, msg any, sender *PID) {
	proc := slf.Registry.get(target)
	if proc != nil {
		proc.Send(target, msg, sender)
	}
}

func (slf *Engine) isLocalMessage(pid *PID) bool {
	return slf.address == pid.Address
}

func (slf *funcReceiver) Receive(ctx *Context) {
	slf.fn(ctx)
}

func newFuncReceiver(fn func(ctx *Context)) Producer {
	return func() Receiver {
		return &funcReceiver{fn: fn}
	}
}
