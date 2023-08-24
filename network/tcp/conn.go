package tcp

import (
	"net"
	"strings"

	"github.com/ZHANG-JIHUI/zephyr/network"
	"github.com/ZHANG-JIHUI/zephyr/third_party/hollywood/actor"
	"github.com/ZHANG-JIHUI/zephyr/tools/log"
	"github.com/panjf2000/gnet/v2"
)

var _ network.Conn = (*tcpConn)(nil)

func newTcpConn(server *actor.PID, gConn gnet.Conn) network.Conn {
	conn := &tcpConn{
		id:     gConn.RemoteAddr().String(),
		addr:   gConn.RemoteAddr(),
		ip:     gConn.RemoteAddr().String(),
		gConn:  gConn,
		server: server,
	}
	if index := strings.LastIndex(conn.ip, ":"); index != -1 {
		conn.ip = conn.ip[0:index]
	}
	return conn
}

type tcpConn struct {
	ctx    *actor.Context
	id     string
	addr   net.Addr
	ip     string
	gConn  gnet.Conn
	server *actor.PID
}

func (slf *tcpConn) PID() *actor.PID {
	if slf.ctx == nil {
		return nil
	}
	return slf.ctx.PID()
}

func (slf *tcpConn) Producer() actor.Producer {
	return func() actor.Receiver {
		return slf
	}
}

func (slf *tcpConn) ID() string {
	return slf.id
}

func (slf *tcpConn) RemoteAddr() net.Addr {
	return slf.addr
}

func (slf *tcpConn) IP() string {
	return slf.ip
}

func (slf *tcpConn) Close() error {
	return slf.gConn.Close()
}

func (slf *tcpConn) Receive(ctx *actor.Context) {
	switch message := ctx.Message().(type) {
	case actor.Started:
		slf.ctx = ctx
	case actor.Stopped:
	default:
		log.Info("tcp conn received message", log.Any("msg", message))
	}
}
