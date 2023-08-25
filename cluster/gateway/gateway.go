package gateway

import (
	"github.com/ZHANG-JIHUI/zephyr/network"
	"github.com/ZHANG-JIHUI/zephyr/network/tcp"
	"github.com/ZHANG-JIHUI/zephyr/third_party/hollywood/actor"
	"github.com/ZHANG-JIHUI/zephyr/tools/concurrent"
	"github.com/ZHANG-JIHUI/zephyr/tools/log"
)

type Gateway struct {
	ctx       *actor.Context
	tcpServer network.Server
	Sessions  *concurrent.Map[*actor.PID, *Session]
}

func NewGateway() *Gateway {
	gate := &Gateway{
		Sessions: concurrent.NewMap[*actor.PID, *Session](),
	}
	tcpServer := tcp.NewServer(":9999")
	tcpServer.RegStartEvent(func(srv network.Server) {
		log.Info("tcp server started", log.String("addr", srv.Addr()), log.Any("pid", srv.PID()))
	})
	tcpServer.RegStopEvent(func(srv network.Server) {
		log.Info("tcp server stopped", log.String("addr", srv.Addr()), log.Any("pid", srv.PID()))
	})
	tcpServer.RegConnectEvent(func(srv network.Server, conn network.Conn) {
		// todo: register pack method
		session := Session{Conn: conn}
		gate.Sessions.Set(conn.PID(), &session)
		log.Info("client connected", log.Any("connection pid", conn.PID()),
			log.Int("online", gate.Sessions.Len()))
	})
	tcpServer.RegDisconnectEvent(func(srv network.Server, conn network.Conn) {
		gate.Sessions.Delete(conn.PID())
		log.Info("client disconnected", log.Any("connection pid", conn.PID()),
			log.Int("online", gate.Sessions.Len()))
	})
	tcpServer.RegReceiveEvent(func(srv network.Server, conn network.Conn, msg []byte, typ int) {
		log.Info("client message", log.Any("pid", conn.PID()), log.ByteString("msg", msg))
		/*
			todo: dispatch client message
		*/
	})
	gate.tcpServer = tcpServer
	return gate
}

func (slf *Gateway) PID() *actor.PID {
	return slf.ctx.PID()
}

func (slf *Gateway) Producer() actor.Producer {
	return func() actor.Receiver {
		return slf
	}
}

func (slf *Gateway) Receive(ctx *actor.Context) {
	switch message := ctx.Message().(type) {
	case actor.Started:
		slf.ctx = ctx
		slf.ctx.SpawnChild(slf.tcpServer.Producer(), "tcp")
	case actor.Stopped:
		slf.ctx.Engine().Poison(slf.tcpServer.PID())
	default:
		log.Info("GATEWAY", log.Any("msg", message))
	}
}
