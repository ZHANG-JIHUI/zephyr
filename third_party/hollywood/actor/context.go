package actor

import (
	"strings"
	"time"

	"github.com/ZHANG-JIHUI/zephyr/tools/concurrent"
	"github.com/ZHANG-JIHUI/zephyr/tools/log"
)

type Context struct {
	engine   *Engine
	pid      *PID
	sender   *PID
	receiver Receiver
	message  any
	parent   *Context
	children *concurrent.Map[string, *PID]
}

func newContext(engine *Engine, pid *PID) *Context {
	return &Context{
		engine:   engine,
		pid:      pid,
		children: concurrent.NewMap[string, *PID](),
	}
}

func (slf *Context) Receiver() Receiver {
	return slf.receiver
}

func (slf *Context) Request(pid *PID, msg any, timeout time.Duration) *Response {
	return slf.engine.Request(pid, msg, timeout)
}

func (slf *Context) Response(msg any) {
	if slf.sender == nil {
		log.Warn("[ACTOR] failed to response because sender is nil", log.Any("pid", slf.pid))
		return
	}
	slf.engine.Send(slf.sender, msg)
}

func (slf *Context) SpawnChild(producer Producer, name string, opts ...Option) *PID {
	options := DefaultOpts(producer)
	options.Name = slf.PID().ID + pidSeparator + name
	for _, opt := range opts {
		opt(&options)
	}
	proc := newProcessor(slf.engine, options)
	proc.ctx.parent = slf
	pid := slf.engine.SpawnProc(proc)
	slf.children.Set(pid.ID, pid)
	return proc.PID()
}

func (slf *Context) SpawnChildFunc(fn func(ctx *Context), name string, opts ...Option) *PID {
	return slf.SpawnChild(newFuncReceiver(fn), name, opts...)
}

// Send 发送消息
func (slf *Context) Send(pid *PID, msg any) {
	slf.engine.SendWithSender(pid, msg, slf.pid)
}

func (slf *Context) SendRepeat(pid *PID, msg any, interval time.Duration) SendRepeater {
	repeater := SendRepeater{
		engine:   slf.engine,
		self:     slf.pid,
		target:   pid.CloneVT(), // 深拷贝，用于新的协程使用
		msg:      msg,
		interval: interval,
		cancel:   make(chan struct{}, 1),
	}
	repeater.start()
	return repeater
}

func (slf *Context) Forward(pid *PID) {
	slf.engine.SendWithSender(pid, slf.message, slf.pid)
}

func (slf *Context) GetPID(name string, tags ...string) *PID {
	if len(tags) > 0 {
		name = name + pidSeparator + strings.Join(tags, pidSeparator)
	}
	proc := slf.engine.Registry.getByID(name)
	if proc != nil {
		return proc.PID()
	}
	return nil
}

func (slf *Context) Parent() *PID {
	if slf.parent != nil {
		return slf.parent.pid
	}
	return nil
}

func (slf *Context) Child(id string) *PID {
	pid, _ := slf.children.Get(id)
	return pid
}

func (slf *Context) Children() []*PID {
	children := make([]*PID, slf.children.Len())
	i := 0
	slf.children.ForEach(func(_ string, pid *PID) {
		children[i] = pid
		i++
	})
	return children
}

func (slf *Context) PID() *PID {
	return slf.pid
}

func (slf *Context) Sender() *PID {
	return slf.sender
}

func (slf *Context) Engine() *Engine {
	return slf.engine
}

func (slf *Context) Message() any {
	return slf.message
}
