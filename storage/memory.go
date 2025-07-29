package storage

import "time"

// MemoryStorage is an in-memory implementation of Storage.
type MemoryStorage struct {
	secrets  map[string]string
	lastSend map[string]time.Time
}

// NewMemoryStorage creates a new MemoryStorage.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		secrets:  make(map[string]string),
		lastSend: make(map[string]time.Time),
	}
}

// GetSecret returns saved secret.
func (m *MemoryStorage) GetSecret(key string) (string, bool) {
	s, ok := m.secrets[key]
	return s, ok
}

// SaveSecret stores secret for key.
func (m *MemoryStorage) SaveSecret(key, secret string) {
	m.secrets[key] = secret
}

// GetLastSend returns last send time.
func (m *MemoryStorage) GetLastSend(key string) (time.Time, bool) {
	t, ok := m.lastSend[key]
	return t, ok
}

// SaveLastSend stores send timestamp.
func (m *MemoryStorage) SaveLastSend(key string, t time.Time) {
	m.lastSend[key] = t
}
