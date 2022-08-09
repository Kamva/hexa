package lg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFnName(t *testing.T) {
	cases := []struct {
		Tag  string
		Fn   any
		Name string
	}{
		{"t1", TestFnName, "TestFnName"},
	}

	for _, c := range cases {
		t.Run(c.Tag, func(t *testing.T) {
			assert.Equal(t, c.Name, FnName(c.Fn))
		})
	}
}
