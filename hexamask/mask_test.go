package hexamask

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type A struct {
	Mask *FieldMask `json:"mask"`
}

func TestFieldMask_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		tag         string
		json        string
		paths       []string
		pathsString string
	}{
		{"t1", `{"mask":null}`, nil, ""},
		{"t2", `{"mask":""}`, []string{}, ""},
		{"t3", `{"mask":"a"}`, []string{"a"}, "a"},
		{"t4", `{"mask":"a,b"}`, []string{"a", "b"}, "a,b"},
		{"t5", `{"mask":"a,b.c"}`, []string{"a", "b.c"}, "a,b.c"},
	}

	for _, test := range tests {
		val := A{}
		assert.Nil(t, json.Unmarshal([]byte(test.json), &val))
		if test.paths != nil {
			assert.Equal(t, test.paths, val.Mask.paths, test.tag)
			assert.Equal(t, test.pathsString, val.Mask.String(), test.tag)
			b, err := json.Marshal(val)
			assert.Nil(t, err, test.tag)
			assert.Equal(t, test.json, string(b), test.tag)
		} else {
			assert.Nil(t, val.Mask, test.tag)
		}
	}
}

func TestFieldMask_IsMasked(t *testing.T) {
	m := FieldMask{paths: []string{"a", "b", "c.d"}}
	assert.False(t, m.IsMasked("d"))
	assert.False(t, m.IsMasked("a.d"))
	assert.True(t, m.IsMasked("a"))
	assert.True(t, m.IsMasked("b"))
	assert.True(t, m.IsMasked("c.d"))
}
