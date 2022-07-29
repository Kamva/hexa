package hexa

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func body(code string, msg string, data any) *HttpRespBody {
	return &HttpRespBody{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

func TestHttpRespBody_MarshalJSONNoMessage(t *testing.T) {
	data := Map{"a": "b"}
	b := body("abc", "", data)
	expected, _ := json.Marshal(Map{"code": "abc", "data": data})

	if m, err := b.MarshalJSON(); assert.NoError(t, err) {
		assert.Equal(t, expected, m)
	}
}

func TestHttpRespBody_MarshalJSONNoData(t *testing.T) {
	b := body("abc", "", nil)
	expected, _ := json.Marshal(Map{"code": "abc"})

	if m, err := b.MarshalJSON(); assert.NoError(t, err) {
		assert.Equal(t, expected, m)
	}
}
