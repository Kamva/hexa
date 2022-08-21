package hlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldToKeyValString(t *testing.T) {
	f := String("a", "b")
	k, v := FieldToKeyVal(f)
	assert.Equal(t, k, "a")
	assert.Equal(t, v, "b")
}

func TestFieldToKeyValInt(t *testing.T) {
	f := Int("a", 1)
	k, v := FieldToKeyVal(f)
	assert.Equal(t, k, "a")
	assert.Equal(t, int64(1), v)
}

func TestFieldToKeyValArr(t *testing.T) {
	val := []string{"a", "b"}
	f := Any("a", val)
	k, v := FieldToKeyVal(f)
	assert.Equal(t, k, "a")
	assert.NotNil(t, v)
}

func TestFieldToKeyValMap(t *testing.T) {
	val := map[string]any{"a": "b"}
	f := Any("a", val)
	k, v := FieldToKeyVal(f)
	assert.Equal(t, k, "a")
	assert.NotNil(t, v)
}
