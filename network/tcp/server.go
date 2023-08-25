package tcp

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/ZHANG-JIHUI/zephyr/network"
	"github.com/ZHANG-JIHUI/zephyr/third_party/hollywood/actor"
	"github.com/ZHANG-JIHUI/zephyr/tools/concurrent"
	"github.com/ZHANG-JIHUI/zephyr/tools/log"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/samber/lo"
)

var _ network.Server = (*server)(nil)
var _ gnet.EventHandler = (*server)(nil)
var _ actor.Receiver = (*server)(nil)

func NewServer(addr string, opts ...Option) network.Server {
	srv := &server{
		Event:       &network.Event{},
		connections: concurrent.NewMap[string, network.Conn](),
	}
	defaultOpts := defaultOptions(addr)
	for _, opt := range opts {
		opt(defaultOpts)
	}
	srv.opts = defaultOpts
	srv.Event.Server = srv
	return srv
}

type server struct {
	*network.Event
	ctx         *actor.Context
	opts        *options
	connections *concurrent.Map[string, network.Conn]
}

func (slf *server) PID() *actor.PID {
	return slf.ctx.PID()
}

func (slf *server) Context() *actor.Context {
	return slf.ctx
}

func (slf *server) Producer() actor.Producer {
	return func() actor.Receiver {
		return slf
	}
}

func (slf *server) Addr() string {
	return slf.opts.addr
}

func (slf *server) Protocol() string {
	return slf.opts.protocol
}

func (slf *server) Start() error {
	go func() {
		if err := gnet.Run(slf, slf.Addr(),
			gnet.WithLogger(log.GetLogger()),
			gnet.WithLogLevel(lo.Ternary(slf.opts.runMode == network.RunModeProd, logging.ErrorLevel, logging.DebugLevel)),
			gnet.WithTicker(true),
			gnet.WithMulticore(true)); err != nil {
			log.Error("gnet server stopped")
			slf.ctx.Engine().Poison(slf.ctx.PID())
		}
	}()
	slf.OnStartEvent()
	return nil
}

func (slf *server) Stop() error {
	slf.OnStopEvent()
	return nil
}

func (slf *server) OnBoot(_ gnet.Engine) (action gnet.Action) {
	return
}

func (slf *server) OnShutdown(eng gnet.Engine) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = eng.Stop(ctx)
}

func (slf *server) OnOpen(gConn gnet.Conn) (out []byte, action gnet.Action) {
	sid := fmt.Sprintf("session_%s", gConn.RemoteAddr().String())
	conn := newTcpConn(slf.ctx.PID(), gConn)
	pid := slf.ctx.SpawnChild(conn.Producer(), sid)
	gConn.SetContext(conn)
	slf.connections.Set(pid.String(), conn)
	slf.OnConnectEvent(conn)
	return
}

func (slf *server) OnClose(gConn gnet.Conn, err error) (action gnet.Action) {
	conn := gConn.Context().(network.Conn)
	pid := conn.PID()
	slf.connections.Delete(pid.String())
	slf.ctx.Engine().Poison(pid)
	slf.OnDisconnectEvent(conn)
	return
}

func (slf *server) OnTraffic(gConn gnet.Conn) (action gnet.Action) {
	conn := gConn.Context().(network.Conn)
	frame, _ := gConn.Next(-1)
	slf.OnReceiveEvent(conn, bytes.Clone(frame), 0)
	return
}

func (slf *server) OnTick() (delay time.Duration, action gnet.Action) {
	delay = time.Second
	return
}

func (slf *server) Receive(ctx *actor.Context) {
	switch message := ctx.Message().(type) {
	case actor.Started:
		slf.ctx = ctx
		if err := slf.Start(); err != nil {
			log.Error("tcp server start err", log.Err(err))
			return
		}
	case actor.Stopped:
		if err := slf.Stop(); err != nil {
			log.Error("tcp server stop err", log.Err(err))
			return
		}
	default:
		fmt.Println(message)
	}
}
