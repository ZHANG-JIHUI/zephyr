package actor

import (
	"context"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type Response struct {
	engine  *Engine
	pid     *PID
	result  chan any
	timeout time.Duration
}

func NewResponse(engine *Engine, timeout time.Duration) *Response {
	return &Response{
		engine:  engine,
		pid:     NewPID(engine.address, "response", strconv.Itoa(rand.Intn(100000))),
		result:  make(chan any, 1),
		timeout: timeout,
	}
}

func (slf *Response) Result() (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), slf.timeout)
	defer func() {
		cancel()
		slf.engine.Registry.Remove(slf.pid)
	}()

	select {
	case resp := <-slf.result:
		return resp, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (slf *Response) Start() {}

func (slf *Response) PID() *PID { return slf.pid }

func (slf *Response) Send(_ *PID, message any, _ *PID) { slf.result <- message }

func (slf *Response) Invoke([]Envelope) {}

func (slf *Response) Shutdown(*sync.WaitGroup) {}
