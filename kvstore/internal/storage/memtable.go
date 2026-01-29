package storage

import "sync"

type MemTable struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewMemTable() *MemTable {
	return &MemTable{data: make(map[string]string)}
}

func (m *MemTable) Put(k, v string) {
	m.mu.Lock()
	m.data[k] = v
	m.mu.Unlock()
}

func (m *MemTable) Get(k string) (string, bool) {
	m.mu.RLock()
	v, ok := m.data[k]
	m.mu.RUnlock()
	return v, ok
}

func (m *MemTable) Snapshot() map[string]string {
	m.mu.RLock()
	cp := make(map[string]string)
	for k, v := range m.data {
		cp[k] = v
	}
	m.mu.RUnlock()
	return cp
}

func (m *MemTable) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

func (m *MemTable) Clear() {
	m.mu.Lock()
	m.data = make(map[string]string)
	m.mu.Unlock()
}
