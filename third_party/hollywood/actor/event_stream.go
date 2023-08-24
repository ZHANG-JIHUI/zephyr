package actor

import (
	"math"
	"math/rand"
	"sync"

	"github.com/ZHANG-JIHUI/zephyr/tools/log"
)

type (
	EventSub    struct{ id uint32 }
	EventFunc   func(event any)
	EventStream struct {
		mux  sync.RWMutex
		subs map[*EventSub]EventFunc
	}
)

func NewEventStream() *EventStream {
	return &EventStream{
		subs: make(map[*EventSub]EventFunc),
	}
}

func (slf *EventStream) Subscribe(fn EventFunc) *EventSub {
	slf.mux.Lock()
	defer slf.mux.Unlock()

	sub := &EventSub{id: uint32(rand.Intn(math.MaxUint32))}
	slf.subs[sub] = fn
	log.Debug("[ACTOR] subscribe event", log.Uint32("id", sub.id), log.Any("fn", fn))

	return sub
}

func (slf *EventStream) Unsubscribe(sub *EventSub) {
	slf.mux.Lock()
	defer slf.mux.Unlock()

	delete(slf.subs, sub)
	log.Debug("[ACTOR] unsubscribe event", log.Uint32("id", sub.id))
}

func (slf *EventStream) Publish(msg any) {
	slf.mux.RLock()
	defer slf.mux.RUnlock()
	for _, fn := range slf.subs {
		go fn(msg)
	}
}

func (slf *EventStream) Len() int {
	slf.mux.RLock()
	defer slf.mux.RUnlock()
	return len(slf.subs)
}
