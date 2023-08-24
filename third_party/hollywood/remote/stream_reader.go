package remote

import (
	"context"
	"errors"

	"github.com/ZHANG-JIHUI/zephyr/third_party/hollywood/actor"
	"github.com/ZHANG-JIHUI/zephyr/tools/log"
)

type streamReader struct {
	DRPCRemoteUnimplementedServer
	remoter      *Remoter
	deserializer Deserializer
}

func newStreamReader(remoter *Remoter) *streamReader {
	return &streamReader{
		remoter:      remoter,
		deserializer: ProtoDeserializer{},
	}
}

// Receive 接收流式消息
func (slf *streamReader) Receive(stream DRPCRemote_ReceiveStream) error {
	defer func() {
		log.Debug("[ACTOR] stream reader terminated")
	}()

	for {
		envelope, err := stream.Recv()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				break
			}
			log.Error("[ACTOR] stream reader receive err", log.Err(err))
			return err
		}
		for _, msg := range envelope.Messages {
			typ := envelope.TypeNames[msg.TypeNameIndex]
			payload, err := slf.deserializer.Deserialize(msg.Data, typ)
			if err != nil {
				log.Error("[ACTOR] stream reader deserialize err", log.Err(err))
				return err
			}
			target := envelope.Targets[msg.TargetIndex]
			var sender *actor.PID
			if len(envelope.Senders) > 0 {
				sender = envelope.Senders[msg.SenderIndex]
			}
			slf.remoter.engine.SendLocal(target, payload, sender)
		}
	}
	return nil
}
