package ringbuffer

import (
	"sync"
	"sync/atomic"
)

type buffer[T any] struct {
	items           []T
	head, tail, mod int64
}

type RingBuffer[T any] struct {
	len     int64
	content *buffer[T]
	mux     sync.Mutex
}

func New[T any](size int64) *RingBuffer[T] {
	return &RingBuffer[T]{
		content: &buffer[T]{
			items: make([]T, size),
			head:  0,
			tail:  0,
			mod:   size,
		},
		len: 0,
	}
}

func (slf *RingBuffer[T]) Push(item T) {
	slf.mux.Lock()
	slf.content.tail = (slf.content.tail + 1) % slf.content.mod
	if slf.content.tail == slf.content.head {
		size := slf.content.mod * 2
		newBuff := make([]T, size)
		for i := int64(0); i < slf.content.mod; i++ {
			idx := (slf.content.tail + i) % slf.content.mod
			newBuff[i] = slf.content.items[idx]
		}
		content := &buffer[T]{
			items: newBuff,
			head:  0,
			tail:  slf.content.mod,
			mod:   size,
		}
		slf.content = content
	}
	atomic.AddInt64(&slf.len, 1)
	slf.content.items[slf.content.tail] = item
	slf.mux.Unlock()
}

func (slf *RingBuffer[T]) Len() int64 {
	return atomic.LoadInt64(&slf.len)
}

func (slf *RingBuffer[T]) Pop() (T, bool) {
	if slf.Len() == 0 {
		var t T
		return t, false
	}
	slf.mux.Lock()
	slf.content.head = (slf.content.head + 1) % slf.content.mod
	item := slf.content.items[slf.content.head]
	var t T
	slf.content.items[slf.content.head] = t
	atomic.AddInt64(&slf.len, -1)
	slf.mux.Unlock()
	return item, true
}

func (slf *RingBuffer[T]) PopN(n int64) ([]T, bool) {
	slf.mux.Lock()
	content := slf.content

	if n >= slf.len {
		n = slf.len
	}
	atomic.AddInt64(&slf.len, -n)

	items := make([]T, n)
	for i := int64(0); i < n; i++ {
		pos := (content.head + 1 + i) % content.mod
		items[i] = content.items[pos]
		var t T
		content.items[pos] = t
	}
	content.head = (content.head + n) % content.mod

	slf.mux.Unlock()
	return items, true
}
