package remote

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/ZHANG-JIHUI/zephyr/third_party/hollywood/actor"
	"github.com/ZHANG-JIHUI/zephyr/tools/log"
	"storj.io/drpc/drpcconn"
)

const (
	connIdleTimeout       = 10 * time.Minute
	streamWriterBatchSize = 1024 * 32
)

func newStreamWriter(engine *actor.Engine, router *actor.PID, address string) actor.Processor {
	return &streamWriter{
		address:    address,
		engine:     engine,
		router:     router,
		pid:        actor.NewPID(engine.Address(), "stream", address),
		inbox:      actor.NewInbox(streamWriterBatchSize),
		serializer: ProtoSerializer{},
	}
}

type streamWriter struct {
	pid        *actor.PID
	engine     *actor.Engine
	address    string
	rawConn    net.Conn
	drpcConn   *drpcconn.Conn
	stream     DRPCRemote_ReceiveStream
	router     *actor.PID
	inbox      *actor.Inbox
	serializer Serializer
}

func (slf *streamWriter) Start() {
	slf.inbox.Start(slf)
	slf.init()
}

func (slf *streamWriter) PID() *actor.PID {
	return slf.pid
}

func (slf *streamWriter) Send(_ *actor.PID, message any, sender *actor.PID) {
	slf.inbox.Send(actor.Envelope{Message: message, Sender: sender})
}

func (slf *streamWriter) Invoke(envelopes []actor.Envelope) {
	var (
		typeNameLookup = make(map[string]int32)
		typeNames      = make([]string, 0)
		senderLookup   = make(map[uint64]int32)
		senders        = make([]*actor.PID, 0)
		targetLookup   = make(map[uint64]int32)
		targets        = make([]*actor.PID, 0)
		messages       = make([]*Message, len(envelopes))
	)

	for i := 0; i < len(envelopes); i++ {
		var (
			stream        = envelopes[i].Message.(*actionDeliver)
			typeNameIndex int32
			senderIndex   int32
			targetIndex   int32
		)
		typeNameIndex, typeNames = slf.lookupTypeName(typeNameLookup, slf.serializer.TypeName(stream.msg), typeNames)
		senderIndex, senders = slf.lookupPIDs(senderLookup, stream.sender, senders)
		targetIndex, targets = slf.lookupPIDs(targetLookup, stream.target, targets)

		data, err := slf.serializer.Serialize(stream.msg)
		if err != nil {
			log.Error("[ACTOR] stream writer serialize err", log.Err(err))
			continue
		}

		messages[i] = &Message{
			Data:          data,
			TypeNameIndex: typeNameIndex,
			SenderIndex:   senderIndex,
			TargetIndex:   targetIndex,
		}
	}

	env := &Packet{
		Senders:   senders,
		Targets:   targets,
		TypeNames: typeNames,
		Messages:  messages,
	}

	if err := slf.stream.Send(env); err != nil {
		if errors.Is(err, io.EOF) {
			_ = slf.rawConn.Close()
			return
		}
		log.Error("[ACTOR] stream writer send err", log.Err(err))
	}
	_ = slf.rawConn.SetDeadline(time.Now().Add(connIdleTimeout))
}

func (slf *streamWriter) Shutdown(wg *sync.WaitGroup) {
	slf.engine.Send(slf.router, &actionTerminate{address: slf.address})
	if slf.stream != nil {
		_ = slf.stream.Close()
	}
	_ = slf.inbox.Stop()
	slf.engine.Registry.Remove(slf.PID())
	if wg != nil {
		wg.Done()
	}
}

func (slf *streamWriter) init() {
	var (
		rawConn net.Conn
		err     error
		delay   = time.Millisecond * 500
	)
	for {
		rawConn, err = net.Dial("tcp", slf.address)
		if err != nil {
			log.Error("[ACTOR] stream writer dial err", log.Err(err))
			time.Sleep(delay)
			continue
		}
		break
	}
	if rawConn == nil {
		slf.Shutdown(nil)
		return
	}

	slf.rawConn = rawConn
	_ = rawConn.SetDeadline(time.Now().Add(connIdleTimeout))

	conn := drpcconn.New(rawConn)
	client := NewDRPCRemoteClient(conn)

	stream, err := client.Receive(context.Background())
	if err != nil {
		log.Error("[ACTOR] stream writer receive error", log.String("address", slf.address), log.Err(err))
	}

	slf.stream = stream
	slf.drpcConn = conn

	log.Info("[ACTOR] stream writer connected", log.String("remote", slf.address))

	go func() {
		<-slf.drpcConn.Closed()
		log.Debug("[ACTOR] stream writer lost connection", log.String("remote", slf.address))
		slf.Shutdown(nil)
	}()
}

func (slf *streamWriter) lookupPIDs(m map[uint64]int32, pid *actor.PID, pids []*actor.PID) (int32, []*actor.PID) {
	if pid == nil {
		return 0, pids
	}
	max := int32(len(m))
	key := pid.LookupKey()
	id, ok := m[key]
	if !ok {
		m[key] = max
		id = max
		pids = append(pids, pid)
	}
	return id, pids

}

func (slf *streamWriter) lookupTypeName(m map[string]int32, name string, types []string) (int32, []string) {
	max := int32(len(m))
	id, ok := m[name]
	if !ok {
		m[name] = max
		id = max
		types = append(types, name)
	}
	return id, types
}
