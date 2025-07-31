package storage

import (
	"testing"
	"time"
)

func TestMemoryStorage_Secret(t *testing.T) {
	m := NewMemoryStorage()
	key := "key"
	s, ok := m.GetSecret(key)
	if ok || s != "" {
		t.Errorf("GetSecret = %v,%v; want \"\",false", s, ok)
	}
	m.SaveSecret(key, "s")
	got, ok := m.GetSecret(key)
	if !ok || got != "s" {
		t.Errorf("GetSecret = %v,%v; want \"s\",true", got, ok)
	}
}

func TestMemoryStorage_LastSend(t *testing.T) {
	m := NewMemoryStorage()
	key := "key"
	t0, ok := m.GetLastSend(key)
	if ok {
		t.Errorf("GetLastSend = %v,true; want _,false", t0)
	}
	now := time.Now()
	m.SaveLastSend(key, now)
	got, ok := m.GetLastSend(key)
	if !ok || !got.Equal(now) {
		t.Errorf("GetLastSend = %v,%v; want %v,true", got, ok, now)
	}
}
