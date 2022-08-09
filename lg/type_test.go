package lg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseType(t *testing.T) {
	cases := []struct {
		Tag           string
		Type          string
		ParsedPackage string
		ParsedType    string
	}{
		{"t1", "", "", ""},
		{"t2", "a", "", "a"},
		{"t3", "a.b", "a", "b"},
		{"t4", "*a.b", "a", "b"},
		{"t5", "[]b", "", "b"},
		{"t5.1", "[]a.b", "a", "b"},
		{"t6", "[]*a.b", "a", "b"},
		{"t7", "map[string]a.b", "a", "b"},
		{"t8", "map[string]*a.b", "a", "b"},
		{"t9", "map[string][]a.b", "a", "b"},
		{"t10", "map[string][]*a.b", "a", "b"},
		{"t11", "map[string]*[]a.b", "a", "b"},
		{"t12", "[]map[string]*[]a.b", "a", "b"},
		{"t13", "[]map[string]*[]b", "", "b"},
	}

	for _, c := range cases {
		t.Run(c.Tag, func(t *testing.T) {
			parsedPackage, parsedType := parseType(c.Type)
			assert.Equal(t, c.ParsedPackage, parsedPackage)
			assert.Equal(t, c.ParsedType, parsedType)
		})
	}
}
func TestSetPackageOnType(t *testing.T) {
	cases := []struct {
		Tag        string
		Type       string
		TargetPkg  string
		ResultType string
	}{
		{"t1", "", "", ""},
		{"t2", "a", "", "a"},
		{"t2", "a.b", "c", "c.b"},
		{"t3", "a.b", "", "b"},
		{"t4", "*a.b", "c", "*c.b"},
		{"t4", "*a.b", "", "*b"},
		{"t5", "[]b", "c", "[]c.b"},
		{"t5.1", "[]a.b", "c", "[]c.b"},
		{"t6", "[]*a.b", "c", "[]*c.b"},
		{"t7", "map[string]a.b", "c", "map[string]c.b"},
		{"t8", "map[string]*a.b", "c", "map[string]*c.b"},
		{"t11", "map[string]*[]a.b", "c", "map[string]*[]c.b"},
		{"t12", "[]map[string]*[]a.b", "c", "[]map[string]*[]c.b"},
		{"t13", "[]map[string]*[]b", "c", "[]map[string]*[]c.b"},
	}

	for _, c := range cases {
		t.Run(c.Tag, func(t *testing.T) {
			assert.Equal(t, c.ResultType, SetPackageOnType(c.TargetPkg, c.Type))
		})
	}
}
