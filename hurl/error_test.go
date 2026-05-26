package hurl

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resp(code int, body string) *http.Response {
	return &http.Response{
		StatusCode: code,
		Status:     http.StatusText(code),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestResponseErr_StatusBoundary(t *testing.T) {
	// 2xx and 3xx are not errors; 4xx/5xx are.
	assert.NoError(t, ResponseErr(resp(http.StatusOK, "")))
	assert.NoError(t, ResponseErr(resp(http.StatusNoContent, "")))
	assert.NoError(t, ResponseErr(resp(http.StatusFound, "")))           // 302 redirect
	assert.NoError(t, ResponseErr(resp(http.StatusMultipleChoices, ""))) // 300

	assert.Error(t, ResponseErr(resp(http.StatusBadRequest, "bad")))           // 400
	assert.Error(t, ResponseErr(resp(http.StatusInternalServerError, "boom"))) // 500
}

func TestResponseErr_CarriesCodeAndBody(t *testing.T) {
	err := ResponseErr(resp(http.StatusNotFound, "missing"))
	require.Error(t, err)

	httpErr, ok := err.(HTTPErr)
	require.True(t, ok)
	assert.Equal(t, http.StatusNotFound, httpErr.Code)
	assert.Equal(t, "missing", httpErr.Body)
}
