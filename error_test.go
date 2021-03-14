package hexa

import (
	"errors"
	"net/http"
	"testing"

	"github.com/kamva/tracer"
	"github.com/stretchr/testify/assert"
)

func TestAsHexaErr(t *testing.T) {
	e,ok:=AsHexaErr(nil)
	assert.Nil(t,e)
	assert.False(t,ok)

	e,ok=AsHexaErr(errors.New("test"))
	assert.Nil(t,e)
	assert.False(t,ok)

	err := NewError(http.StatusBadRequest, "a", nil)

	hexaErr, ok := AsHexaErr(err)
	assert.NotNil(t, hexaErr)
	assert.True(t, ok)

	hexaErr, ok = AsHexaErr(tracer.Trace(err))
	assert.NotNil(t, hexaErr)
	assert.True(t, ok)
}
