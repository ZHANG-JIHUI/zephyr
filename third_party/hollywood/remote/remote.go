package remote

import (
	"context"
	"net"

	"github.com/ZHANG-JIHUI/zephyr/third_party/hollywood/actor"
	"github.com/ZHANG-JIHUI/zephyr/tools/log"
	"storj.io/drpc/drpcmux"
	"storj.io/drpc/drpcserver"
)

type Remoter struct {
	engine  *actor.Engine // actor引擎
	reader  *streamReader // reader
	router  *actor.PID    // 路由器
	address string        // 监听地址
}

func NewRemoter(engine *actor.Engine, address string) *Remoter {
	remoter := &Remoter{
		engine:  engine,
		address: address,
	}
	remoter.reader = newStreamReader(remoter)
	return remoter
}

func (slf *Remoter) Start() {
	ln, err := net.Listen("tcp", slf.address)
	if err != nil {
		log.Fatal("[ACTOR] remoter listen err", log.Err(err))
	}
	mux := drpcmux.New()
	_ = DRPCRegisterRemote(mux, slf.reader)
	server := drpcserver.New(mux)
	slf.router = slf.engine.Spawn(newStreamRouter(slf.engine), "router", actor.WithInboxSize(1024*1024))

	log.Info("[ACTOR] remoter started", log.String("address", slf.address))

	go func() {
		_ = server.Serve(context.Background(), ln)
	}()
}

func (slf *Remoter) Send(target *actor.PID, msg any, sender *actor.PID) {
	slf.engine.Send(slf.router, &actionDeliver{
		target: target,
		msg:    msg,
		sender: sender,
	})
}

func (slf *Remoter) Address() string {
	return slf.address
}
