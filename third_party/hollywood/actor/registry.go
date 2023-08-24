package actor

import (
	"sync"

	"github.com/ZHANG-JIHUI/zephyr/tools/log"
)

const LocalLookupAddr = "local"

type Registry struct {
	mux    sync.RWMutex
	lookup map[string]Processor
	engine *Engine
}

func newRegistry(engine *Engine) *Registry {
	return &Registry{
		lookup: make(map[string]Processor, 1024),
		engine: engine,
	}
}

func (slf *Registry) add(proc Processor) {
	slf.mux.Lock()
	defer slf.mux.Unlock()
	id := proc.PID().ID
	if _, ok := slf.lookup[id]; ok {
		log.Warn("[ACTOR] processor already registered", log.Any("pid", proc.PID()))
		return
	}
	slf.lookup[id] = proc
}

func (slf *Registry) Remove(pid *PID) {
	slf.mux.Lock()
	defer slf.mux.Unlock()
	delete(slf.lookup, pid.ID)
}

func (slf *Registry) get(pid *PID) Processor {
	slf.mux.RLock()
	defer slf.mux.RUnlock()
	if proc, ok := slf.lookup[pid.ID]; ok {
		return proc
	}
	return slf.engine.deadLetter
}

func (slf *Registry) getByID(id string) Processor {
	slf.mux.RLock()
	slf.mux.RUnlock()
	return slf.lookup[id]
}
