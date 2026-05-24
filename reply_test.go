package hexa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReply(t *testing.T) {
	r := NewReply(200, "lib.ok")
	assert.Equal(t, 200, r.HTTPStatus())
	assert.Equal(t, "lib.ok", r.ID())
	assert.Nil(t, r.Data())

	// Setters return a modified copy and leave the original untouched.
	r2 := r.SetHTTPStatus(201).SetData(Map{"a": 1})
	assert.Equal(t, 201, r2.HTTPStatus())
	assert.Equal(t, Map{"a": 1}, r2.Data())

	assert.Equal(t, 200, r.HTTPStatus())
	assert.Nil(t, r.Data())
}
