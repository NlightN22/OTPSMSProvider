package storage

import "time"

// Storage defines methods to persist secrets and timestamps.
type Storage interface {
	GetSecret(key string) (string, bool)
	SaveSecret(key, secret string)
	GetLastSend(key string) (time.Time, bool)
	SaveLastSend(key string, t time.Time)
}
