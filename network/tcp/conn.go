package tcp

import (
	"net"
	"strings"
	"sync"
	"time"

	"github.com/ZHANG-JIHUI/zephyr/network"
	"github.com/ZHANG-JIHUI/zephyr/third_party/hollywood/actor"
	"github.com/ZHANG-JIHUI/zephyr/tools/concurrent"
	"github.com/panjf2000/gnet/v2"
)

var _ network.Conn = (*tcpConn)(nil)

func newTcpConn(server *actor.PID, gn gnet.Conn) network.Conn {
	conn := &tcpConn{
		id:     gn.RemoteAddr().String(),
		addr:   gn.RemoteAddr(),
		ip:     gn.RemoteAddr().String(),
		gn:     gn,
		server: server,
		close:  make(chan struct{}, 1),
	}
	if index := strings.LastIndex(conn.ip, ":"); index != -1 {
		conn.ip = conn.ip[0:index]
	}
	return conn
}

type tcpConn struct {
	ctx          *actor.Context
	id           string
	addr         net.Addr
	ip           string
	gn           gnet.Conn
	server       *actor.PID
	close        chan struct{}
	mutex        sync.Mutex
	envelopes    []*network.Envelope
	envelopePool *concurrent.Pool[*network.Envelope]
}

func (slf *tcpConn) PID() *actor.PID {
	if slf.ctx == nil {
		return nil
	}
	return slf.ctx.PID()
}

func (slf *tcpConn) Context() *actor.Context {
	return slf.ctx
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
	_ = slf.gn.Close()
	if slf.envelopePool != nil {
		slf.envelopePool.Close()
	}
	slf.envelopePool = nil
	slf.envelopes = nil
	return nil
}

func (slf *tcpConn) Write(data []byte) {
	if slf.envelopePool == nil {
		return
	}
	envelope := slf.envelopePool.Get()
	envelope.Data = data
	slf.mutex.Lock()
	slf.envelopes = append(slf.envelopes, envelope)
	slf.mutex.Unlock()
}

func (slf *tcpConn) Receive(ctx *actor.Context) {
	switch message := ctx.Message().(type) {
	case actor.Started:
		slf.ctx = ctx
		wg := sync.WaitGroup{}
		wg.Add(1)
		go slf.writeLoop(&wg)
		wg.Wait()
	case actor.Stopped:
		_ = slf.Close()
	case *network.NetPacketMessage:
		// TODO：pack message to byte array
		slf.Write(message.Data)
	default:
	}
}

func (slf *tcpConn) writeLoop(wg *sync.WaitGroup) {
	slf.envelopePool = concurrent.NewPool[*network.Envelope](1024*10,
		func() *network.Envelope { return &network.Envelope{} },
		func(envelope *network.Envelope) { envelope.Data = nil },
	)
	defer func() {
		if err := recover(); err != nil {
			slf.ctx.Engine().Poison(slf.PID())
		}
	}()
	wg.Done()

	for {
		slf.mutex.Lock()
		if slf.envelopePool == nil {
			return
		}
		if len(slf.envelopes) == 0 {
			slf.mutex.Unlock()
			time.Sleep(50 * time.Millisecond)
			continue
		}
		envelopes := slf.envelopes[0:]
		slf.envelopes = slf.envelopes[0:0]
		slf.mutex.Unlock()
		for i := 0; i < len(envelopes); i++ {
			env := envelopes[i]
			err := slf.gn.AsyncWrite(env.Data, nil)
			slf.envelopePool.Release(env)
			if err != nil {
				panic(err)
			}
		}
	}
}
