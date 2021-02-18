package store

import (
	"fmt"
	"sync"
)

var _ DataStore = (*Memory)(nil)

// Memory implements the DataStore interface and provides a local
// in-memory, thread-safe key/value store.
type Memory struct {
	ErrorOnDuplicate bool // Should we return an error on a duplicate key?

	mu   sync.RWMutex
	data map[string]string
}

// NewMemory returns a new Memory data store adapter.
func NewMemory() *Memory {
	return &Memory{
		mu:   sync.RWMutex{},
		data: make(map[string]string),
	}
}

// Set adds or updates a value.
// Returns an error if ErrorOnDuplicate is true.
func (m *Memory) Set(key string, value string) error {
	// Note explicitly calling unlock vs deferring it in order to
	// unlock as fast as possible.
	m.mu.Lock()
	if m.ErrorOnDuplicate {
		if _, ok := m.data[key]; ok {
			m.mu.Unlock()
			return fmt.Errorf("Error setting key %s: %w", truncate(key, 32), ErrKeyExists)
		}
	}
	m.data[key] = value
	m.mu.Unlock()
	return nil
}

// Get attempts to get a value by the given key.
// Returns an error if the key is non-existent.
func (m *Memory) Get(key string) (*string, error) {
	// Note explicitly calling unlock vs deferring it in order to
	// unlock as fast as possible.
	m.mu.RLock()
	val, ok := m.data[key]
	m.mu.RUnlock()

	if ok {
		return &val, nil
	}

	return nil, fmt.Errorf("Error: %w", ErrKeyNotFound)
}

// Del removes all provided keys provided they exist.
// Returns the number of keys removed. Does nothing if a key
// does not exist.
func (m *Memory) Del(keys []string) int {
	count := 0
	m.mu.Lock()
	for _, key := range keys {
		if _, ok := m.data[key]; ok {
			count++
			delete(m.data, key)
		}
	}
	m.mu.Unlock()
	return count
}
