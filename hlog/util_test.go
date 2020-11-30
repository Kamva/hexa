package hlog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFieldToKeyValString(t *testing.T) {
	f := String("a", "b")
	k, v := fieldToKeyVal(f)
	assert.Equal(t, k, "a")
	assert.Equal(t, v, "b")
}

func TestFieldToKeyValInt(t *testing.T) {
	f := Int("a", 1)
	k, v := fieldToKeyVal(f)
	assert.Equal(t, k, "a")
	assert.Equal(t, int64(1), v)
}

func TestFieldToKeyValArr(t *testing.T) {
	val := []string{"a", "b"}
	f := Any("a", val)
	k, v := fieldToKeyVal(f)
	assert.Equal(t, k, "a")
	assert.NotNil(t, v)
}

func TestFieldToKeyValMap(t *testing.T) {
	val := map[string]interface{}{"a": "b"}
	f := Any("a", val)
	k, v := fieldToKeyVal(f)
	assert.Equal(t, k, "a")
	assert.NotNil(t, v)
}
