package hurl

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedactHeaders(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "Bearer secret-token")
	h.Add("Cookie", "session=abc")
	h.Set("Content-Type", "application/json")

	red := redactHeaders(h)

	// Sensitive headers are masked in the copy.
	assert.Equal(t, "REDACTED", red.Get("Authorization"))
	assert.Equal(t, "REDACTED", red.Get("Cookie"))
	// Non-sensitive headers are preserved.
	assert.Equal(t, "application/json", red.Get("Content-Type"))
	// The original header is untouched.
	assert.Equal(t, "Bearer secret-token", h.Get("Authorization"))

	// nil-safe.
	assert.Nil(t, redactHeaders(nil))
}

func Test_isValidUrl(t *testing.T) {
	urls := []struct {
		tag     string
		val     string
		isValid bool
	}{
		{"test 1", "https://abc.com", true},
		{"test 2", "https://", false},
		{"test 3", "abc", false},
		{"test 4", "abc.com", false},
		{"test 5", "www.abc.com/a/b", false},
		{"test 5", "https://abc.com/a/b", true},
		{"test 6", "/a/b", false},
		{"test 7", "?abc/d", false},
	}

	for _, u := range urls {
		t.Run(u.tag, func(t *testing.T) {
			assert.Equal(t, u.isValid, isValidURL(u.val))
		})
	}
}
