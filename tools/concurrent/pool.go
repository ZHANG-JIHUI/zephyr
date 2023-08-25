package concurrent

import (
	"fmt"
	"sync"

	"github.com/ZHANG-JIHUI/zephyr/tools/log"
)

type Pool[T any] struct {
	mutex     sync.Mutex
	buffers   []T
	size      int
	generator func() T
	releaser  func(data T)
	warn      int
}

func NewPool[T any](size int, generator func() T, releaser func(data T)) *Pool[T] {
	pool := &Pool[T]{
		size:      size,
		generator: generator,
		releaser:  releaser,
	}
	for i := 0; i < size; i++ {
		pool.put(generator())
	}
	return pool
}

func (slf *Pool[T]) Get() T {
	slf.mutex.Lock()
	if len(slf.buffers) > 0 {
		data := slf.buffers[0]
		slf.buffers = slf.buffers[1:]
		slf.mutex.Unlock()
		return data
	}
	slf.mutex.Unlock()
	slf.warn++
	if slf.warn >= slf.size/10 || slf.warn >= 1000 {
		log.Warn("Pool", log.String("Get", "the number of buffer members is insufficient, consider whether it is due to unreleased or inappropriate buffer size"), log.String("Status", fmt.Sprintf("%d/%d", slf.size+slf.warn, slf.size)))
		slf.warn = 0
	}
	return slf.generator()
}

func (slf *Pool[T]) Release(data T) {
	slf.releaser(data)
	slf.put(data)
}

func (slf *Pool[T]) Close() {
	slf.mutex.Lock()
	slf.buffers = nil
	slf.size = 0
	slf.generator = nil
	slf.releaser = nil
	slf.warn = 0
	slf.mutex.Unlock()
}

func (slf *Pool[T]) IsClosed() bool {
	return slf.generator == nil
}

func (slf *Pool[T]) put(data T) {
	slf.mutex.Lock()
	if len(slf.buffers) >= slf.size {
		slf.mutex.Unlock()
		return
	}
	slf.buffers = append(slf.buffers, data)
	slf.mutex.Unlock()
}
