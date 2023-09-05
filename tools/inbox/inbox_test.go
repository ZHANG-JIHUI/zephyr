package inbox

import (
	"fmt"
	"testing"
	"time"

	"github.com/ZHANG-JIHUI/zephyr/tools/log"
)

type processor struct {
	inbox *Inbox
}

type Message struct {
	ID   int64
	Data []byte
}

func (slf *processor) Invoke(envelopes []Envelope) {
	for _, envelope := range envelopes {
		switch msg := envelope.Msg.(type) {
		case string:
			log.Info("processor.Invoke", log.String("msg", msg))
		case Message:
			log.Info("processor.Invoke", log.Int64("id", msg.ID), log.ByteString("data", msg.Data))
		}
	}
}

func TestInbox(t *testing.T) {
	proc := &processor{inbox: NewInbox[int32](10)}
	proc.inbox.Start(proc)

	for i := 0; i < 1000; i++ {
		proc.inbox.Push(Envelope{Msg: fmt.Sprintf("hello %d", i)})
	}

	time.Sleep(10 * time.Second)

	_ = proc.inbox.Stop()
}
