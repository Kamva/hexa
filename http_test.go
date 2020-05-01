package hexa

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHttpRespBody_MarshalJSONNoMessage(t *testing.T) {
	data := Map{"a": "b"}
	b := NewBody("abc", "", data)
	expected, _ := json.Marshal(Map{"code": "abc", "data": data})

	if m, err := b.MarshalJSON(); assert.NoError(t, err) {
		assert.Equal(t, expected, m)
	}
}

func TestHttpRespBody_MarshalJSONNoData(t *testing.T) {
	b := NewBody("abc", "", nil)
	expected, _ := json.Marshal(Map{"code": "abc"})

	if m, err := b.MarshalJSON(); assert.NoError(t, err) {
		assert.Equal(t, expected, m)
	}
}
