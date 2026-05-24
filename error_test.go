package hexa

import (
	"errors"
	"net/http"
	"testing"

	"github.com/kamva/tracer"
	"github.com/stretchr/testify/assert"
)

func TestAsHexaErr(t *testing.T) {
	e := AsHexaErr(nil)
	assert.Nil(t, e)

	e = AsHexaErr(errors.New("test"))
	assert.Nil(t, e)

	err := NewError(http.StatusBadRequest, "a")

	hexaErr := AsHexaErr(err)
	assert.NotNil(t, hexaErr)

	hexaErr = AsHexaErr(tracer.Trace(err))
	assert.NotNil(t, hexaErr)
}

func TestDefaultError_Unwrap(t *testing.T) {
	sentinel := errors.New("sentinel")
	err := NewError(http.StatusBadRequest, "lib.x").SetError(sentinel)

	// errors.Is must reach the wrapped internal error through Unwrap.
	assert.True(t, errors.Is(err, sentinel))

	// A non-matching target stays false.
	assert.False(t, errors.Is(err, errors.New("other")))
}
