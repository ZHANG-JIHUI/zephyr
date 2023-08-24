package remote

import (
	"github.com/ZHANG-JIHUI/zephyr/third_party/hollywood/actor"
	"github.com/ZHANG-JIHUI/zephyr/tools/log"
)

type (
	actionDeliver struct {
		sender *actor.PID
		target *actor.PID
		msg    any
	}
	actionTerminate struct {
		address string
	}
)

type streamRouter struct {
	engine  *actor.Engine
	streams map[string]*actor.PID
	pid     *actor.PID
}

func newStreamRouter(engine *actor.Engine) actor.Producer {
	return func() actor.Receiver {
		return &streamRouter{
			engine:  engine,
			streams: make(map[string]*actor.PID),
		}
	}
}

func (slf *streamRouter) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		slf.pid = ctx.PID()
	case *actionDeliver:
		slf.deliverStream(msg)
	case *actionTerminate:
		slf.terminateStream(msg)
	}
}

func (slf *streamRouter) deliverStream(payload *actionDeliver) {

	address := payload.target.Address
	pid, exist := slf.streams[address]
	if !exist {
		pid = slf.engine.SpawnProc(newStreamWriter(slf.engine, slf.pid, address))
		slf.streams[address] = pid
		log.Debug("[ACTOR] new stream router", log.Any("pid", pid))
	}
	slf.engine.Send(pid, payload)
}

func (slf *streamRouter) terminateStream(payload *actionTerminate) {
	pid := slf.streams[payload.address]
	delete(slf.streams, payload.address)
	log.Debug("[ACTOR] stream router terminating stream",
		log.String("remote", payload.address), log.Any("pid", pid))
}
