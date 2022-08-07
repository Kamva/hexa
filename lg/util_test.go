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
	}

	for _, c := range cases {
		t.Run(c.Tag, func(t *testing.T) {
			parsedPackage, parsedType := parseType(c.Type)
			assert.Equal(t, c.ParsedPackage, parsedPackage)
			assert.Equal(t, c.ParsedType, parsedType)
		})
	}
}
