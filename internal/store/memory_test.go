package store

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemory_Set(t *testing.T) {
	kvPairs := [][]string{
		{"1", "1"},
	}

	m := NewMemory()

	for i := 0; i < len(kvPairs); i++ {
		k := kvPairs[i][0]
		v := kvPairs[i][1]

		m.Set(k, v)

		assert.Equal(t, v, m.data[k])
	}
}

func TestMemory_Get(t *testing.T) {
	s := strings.Repeat("abc", 256)

	kvPairs := [][]string{
		{"1", "1"},
		{s, s},
	}

	m := NewMemory()

	for i := 0; i < len(kvPairs); i++ {
		k := kvPairs[i][0]
		v := kvPairs[i][1]

		m.data[k] = v

		r, err := m.Get(k)

		assert.Nil(t, err)
		assert.NotNil(t, r)

		if r != nil {
			assert.Equal(t, v, *r)
		}
	}
}

func TestMemory_Del(t *testing.T) {
	s := strings.Repeat("abc", 256)

	kvPairs := [][]string{
		{"1", "1"},
		{"hello", "world"},
		{s, s},
	}

	keys := []string{}

	m := NewMemory()

	for i := 0; i < len(kvPairs); i++ {
		k := kvPairs[i][0]
		v := kvPairs[i][1]

		m.data[k] = v
		keys = append(keys, k)
	}

	count := m.Del(keys)
	assert.Equal(t, len(keys), count)
}
