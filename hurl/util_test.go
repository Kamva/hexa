package hurl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
