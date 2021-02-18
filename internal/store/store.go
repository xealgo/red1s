package store

import (
	"github.com/xealgo/red1s"
)

// Error constants
const (
	ErrKeyExists   = red1s.Error("Duplicate key found")
	ErrKeyNotFound = red1s.Error("Key not found")
)

// DataStore used to abstract the underlying storage mechanism.
type DataStore interface {
	Set(key string, value string) error
	Get(key string) (*string, error)
	Del(keys []string) int
}

func truncate(s string, max int) string {
	if max > len(s) {
		return s
	}
	return s[:max] + ".."
}
