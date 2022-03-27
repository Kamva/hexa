package hexa

import (
	"testing"
	"time"

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

func TestStoreAtomic(t *testing.T) {
	origS := newStore()

	origS.Atomic(func(s Store) {
		// Store inside the atomic function should panic when we call to the atomic.
		assert.Panics(t, func() {
			s.Atomic(func(s Store) {})
		})
		s.Set("a", "123")

		// It shouldn't let us to set data:
		go func() {
			origS.Set("a", "abc")
		}()

		// Wait to make sure it had enough time to set the value using origS.
		time.Sleep(time.Millisecond * 10)
		assert.Equal(t, "123", s.Get("a").(string))
	})
}
