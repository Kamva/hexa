package hexa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreImpl(t *testing.T) {
	s := newStore()
	assert.Nil(t, s.Get("abc"))

	s.(*atomicStore).Set("abc", "cde")
	assert.Equal(t, "cde", s.Get("abc"))

	s.(*atomicStore).Set("abc", "123")
	assert.Equal(t, "123", s.Get("abc"))
}

func TestAtomicStore_SetIfNotExist(t *testing.T) {
	s := newStore()
	s.Set("a", "abc")
	s.SetIfNotExist("a", func() interface{} {
		return "def"
	})
	s.SetIfNotExist("b", func() interface{} {
		return "123"
	})

	assert.Equal(t, "abc", s.Get("a").(string))
	assert.Equal(t, "123", s.Get("b").(string))
}

func BenchmarkSetIfNotExist(b *testing.B) {
	s := newStore()
	for n := 0; n < b.N; n++ {
		val := s.SetIfNotExist("a", func() interface{} {
			return "abc"
		})
		_ = val
	}
}
