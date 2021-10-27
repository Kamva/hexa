package hexa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreImpl(t *testing.T) {
	s := newStore()
	assert.Nil(t, s.Get("abc"))

	s.(*storeImpl).Set("abc","cde")
	assert.Equal(t, "cde", s.Get("abc"))

	s.(*storeImpl).Set("abc","123")
	assert.Equal(t, "123", s.Get("abc"))
}
