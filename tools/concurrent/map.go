package concurrent

import (
	"encoding/json"
	"sync"
)

func NewMap[K comparable, V any](options ...MapOption[K, V]) *Map[K, V] {
	m := &Map[K, V]{
		data: make(map[K]V),
	}
	for _, option := range options {
		option(m)
	}
	return m
}

type Map[K comparable, V any] struct {
	mux  sync.RWMutex
	data map[K]V
}

func (slf *Map[K, V]) Set(k K, v V) {
	slf.mux.Lock()
	defer slf.mux.Unlock()
	slf.data[k] = v
}

func (slf *Map[K, V]) Get(k K) (V, bool) {
	slf.mux.RLock()
	defer slf.mux.RUnlock()
	val, ok := slf.data[k]
	return val, ok
}

func (slf *Map[K, V]) Delete(k K) {
	slf.mux.Lock()
	defer slf.mux.Unlock()
	delete(slf.data, k)
}

func (slf *Map[K, V]) Exist(key K) bool {
	slf.mux.RLock()
	defer slf.mux.RUnlock()
	_, exist := slf.data[key]
	return exist
}

func (slf *Map[K, V]) Len() int {
	slf.mux.RLock()
	defer slf.mux.RUnlock()
	return len(slf.data)
}

func (slf *Map[K, V]) ForEach(f func(K, V)) {
	slf.mux.RLock()
	defer slf.mux.RUnlock()
	for key, val := range slf.data {
		f(key, val)
	}
}

func (slf *Map[K, V]) Keys() []K {
	slf.mux.RLock()
	defer slf.mux.RUnlock()
	keys := make([]K, 0, len(slf.data))
	for key := range slf.data {
		keys = append(keys, key)
	}
	return keys
}

func (slf *Map[K, V]) Values() []V {
	slf.mux.RLock()
	defer slf.mux.RUnlock()
	values := make([]V, 0, len(slf.data))
	for _, value := range slf.data {
		values = append(values, value)
	}
	return values
}

func (slf *Map[K, V]) Map() map[K]V {
	slf.mux.RLock()
	defer slf.mux.RUnlock()
	m := make(map[K]V)
	for k, v := range slf.data {
		m[k] = v
	}
	return m
}

func (slf *Map[K, V]) MarshalJSON() ([]byte, error) {
	m := slf.Map()
	return json.Marshal(m)
}

func (slf *Map[K, V]) UnmarshalJSON(bytes []byte) error {
	m := make(map[K]V)
	slf.mux.Lock()
	slf.mux.Unlock()
	if err := json.Unmarshal(bytes, &m); err != nil {
		return err
	}
	slf.data = m
	return nil
}
