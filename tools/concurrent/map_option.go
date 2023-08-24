package concurrent

type MapOption[K comparable, V any] func(*Map[K, V])

func WithMapSource[K comparable, V any](source map[K]V) MapOption[K, V] {
	return func(m *Map[K, V]) {
		m.data = source
	}
}
