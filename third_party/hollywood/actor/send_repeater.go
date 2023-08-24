package actor

import "time"

type SendRepeater struct {
	engine   *Engine
	self     *PID
	target   *PID
	msg      any
	interval time.Duration
	cancel   chan struct{}
}

func (slf *SendRepeater) start() {
	ticker := time.NewTicker(slf.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				slf.engine.SendWithSender(slf.target, slf.msg, slf.self)
			case <-slf.cancel:
				ticker.Stop()
				return
			}
		}
	}()
}

func (slf *SendRepeater) stop() {
	close(slf.cancel)
}
