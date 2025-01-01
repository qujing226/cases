package syncMap

import "sync"

type Map[k comparable, v any] struct {
	mu sync.RWMutex
	m  map[k]v
}

func NewMap[k comparable, v any]() *Map[k, v] {
	return &Map[k, v]{
		m: make(map[k]v),
	}
}

func (m *Map[k, v]) Load(key k) (v, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.m[key]
	if ok {
		return v, true
	}
	return v, false
}

func (m *Map[k, v]) Store(key k, val v) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[key] = val
}
