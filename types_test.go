package hexa

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecret_Masks(t *testing.T) {
	s := Secret("super-secret")

	assert.Equal(t, "****", s.String())
	assert.Equal(t, "****", fmt.Sprintf("%s", s))

	b, err := s.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `"****"`, string(b))
}

func TestExtendBytesMap(t *testing.T) {
	dest := map[string][]byte{"a": []byte("1")}
	src := map[string][]byte{"a": []byte("2"), "b": []byte("3")}

	// Without overwrite, existing keys are kept.
	extendBytesMap(dest, src, false)
	assert.Equal(t, "1", string(dest["a"]))
	assert.Equal(t, "3", string(dest["b"]))

	// With overwrite, existing keys are replaced.
	extendBytesMap(dest, map[string][]byte{"a": []byte("9")}, true)
	assert.Equal(t, "9", string(dest["a"]))
}
