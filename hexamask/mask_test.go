package hexamask

import (
	"encoding/json"
	"reflect"
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

func TestFieldMask_PathIsMasked(t *testing.T) {
	m := FieldMask{paths: []string{"a", "b", "c.d"}}
	assert.False(t, m.PathIsMasked("d"))
	assert.False(t, m.PathIsMasked("a.d"))
	assert.True(t, m.PathIsMasked("a"))
	assert.True(t, m.PathIsMasked("b"))
	assert.True(t, m.PathIsMasked("c.d"))
}

func TestFieldMask_pathOf(t *testing.T) {
	tests := []struct {
		tag      string
		fieldTag string
		path     string
	}{
		{"t1", `mask:"mask_val"`, "mask_val"},
		{"t2", `json:"json_val" mask:"mask_val"`, "mask_val"},
		{"t3", `json:"json_val"`, "json_val"},
		{"t4", ``, ""},
	}

	m := FieldMask{}
	for _, test := range tests {
		assert.Equal(t, test.path, m.pathOf(reflect.StructTag(test.fieldTag)), test.tag)
	}
}

func TestFieldMask_Mask(t *testing.T) {
	type B struct {
		Letter string `json:"letter"`
		Val    string `json:"val"`
	}
	type A struct {
		Name   string `json:"name"`
		Age    *int   `json:"age"`
		Salary int    `json:"salary"`
		B      B      `json:"b"`
		B2     B
		Mask   *FieldMask `json:"mask"`
	}

	inp := A{
		Name: "abc",
		Age:  nil,
		B:    B{Val: "salam"},
		Mask: &FieldMask{paths: []string{"c.d", "name", "age", "b", "b.val"}},
	}
	inp.Mask.Mask(&inp)
	assert.Equal(t, inp.Mask.maskedFields, []interface{}{
		&inp.Name,
		&inp.Age,
		&inp.B,
		&inp.B.Val,
	})
}

func TestFieldMask_IsMasked(t *testing.T) {
	type B struct {
		Letter string `json:"letter"`
		Val    string `json:"val"`
	}
	type A struct {
		Name   string `json:"name"`
		Age    *int   `json:"age"`
		Salary int    `json:"salary"`
		B      B      `json:"b"`
		B2     B
	}

	inp := A{
		Name: "abc",
		Age:  nil,
		B:    B{Val: "salam"},
	}

	m := FieldMask{maskedFields: []interface{}{
		&inp.Name,
		&inp.Age,
		&inp.B,
		&inp.B.Val,
	}}
	assert.True(t, m.IsMasked(&inp.Name))
	assert.True(t, m.IsMasked(&inp.Age))
	assert.False(t, m.IsMasked(&inp.Salary))
	assert.True(t, m.IsMasked(&inp.B))
	assert.False(t, m.IsMasked(&inp.B.Letter))
	assert.True(t, m.IsMasked(&inp.B.Val))
	assert.False(t, m.IsMasked(&inp.B2))
	assert.False(t, m.IsMasked(&inp.B2.Letter))
	assert.False(t, m.IsMasked(&inp.B2.Val))
}

func TestFieldMask_MaskStrut(t *testing.T) {
	type B struct {
		Letter string `json:"letter"`
		Val    string `json:"val"`
	}
	type A struct {
		Name   string `json:"name"`
		Age    *int   `json:"age"`
		Salary int    `json:"salary"`
		B      B      `json:"b"`
		B2     B
		Mask   *FieldMask `json:"mask"`
	}

	inp := A{
		Name: "abc",
		Age:  nil,
		B:    B{Val: "salam"},
		Mask: &FieldMask{paths: []string{"c.d", "name", "age", "b", "b.val"}},
	}
	inp.Mask.Mask(&inp)

	assert.True(t, inp.Mask.IsMasked(&inp.Name))
	assert.True(t, inp.Mask.IsMasked(&inp.Age))
	assert.False(t, inp.Mask.IsMasked(&inp.Salary))
	assert.True(t, inp.Mask.IsMasked(&inp.B))
	assert.False(t, inp.Mask.IsMasked(&inp.B.Letter))
	assert.True(t, inp.Mask.IsMasked(&inp.B.Val))
	assert.False(t, inp.Mask.IsMasked(&inp.B2))
	assert.False(t, inp.Mask.IsMasked(&inp.B2.Letter))
	assert.False(t, inp.Mask.IsMasked(&inp.B2.Val))
}
