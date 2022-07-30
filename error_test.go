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
