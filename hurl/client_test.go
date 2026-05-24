package hurl

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/kamva/hexa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURL_BuildsFromBase(t *testing.T) {
	u, err := NewURL("http://a.com/api")
	require.NoError(t, err)

	cases := map[string]string{
		"users":             "http://a.com/api/users", // relative joins onto base path
		"/v2/x":             "http://a.com/v2/x",      // leading slash replaces base path
		"http://b.com/full": "http://b.com/full",      // absolute URL is used as-is
	}
	for in, want := range cases {
		got, err := u.URL(in)
		require.NoError(t, err, in)
		assert.Equal(t, want, got.String(), in)
	}
}

func TestURL_EmptyBaseRequiresAbsolute(t *testing.T) {
	u, err := NewURL("")
	require.NoError(t, err)

	// A relative path without a base is an error, not a panic.
	_, err = u.URL("relative/path")
	assert.Error(t, err)

	// An absolute URL still works without a base.
	got, err := u.URL("https://x.com/y")
	require.NoError(t, err)
	assert.Equal(t, "https://x.com/y", got.String())
}

func TestURLQueryParams(t *testing.T) {
	u, err := NewURL("http://a.com/api")
	require.NoError(t, err)

	got, err := u.URL("search", URLQueryParams(hexa.Map{"q": "go", "n": 3}))
	require.NoError(t, err)
	assert.Equal(t, "go", got.Query().Get("q"))
	assert.Equal(t, "3", got.Query().Get("n"))
}

func TestRequestOptions_SetHeaders(t *testing.T) {
	tests := []struct {
		name   string
		opt    RequestOption
		header string
		want   string
	}{
		{"basic auth", BasicAuth("user", "pass"), "Authorization", "Basic dXNlcjpwYXNz"},
		{"bearer", BearerToken("tok"), "Authorization", "Bearer tok"},
		{"custom type", AuthorizationToken("Token", "abc"), "Authorization", "Token abc"},
		{"empty type", AuthenticateHeader("X-Auth", "", "raw"), "X-Auth", "raw"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://a.com", nil)
			require.NoError(t, tc.opt(req))
			assert.Equal(t, tc.want, req.Header.Get(tc.header))
		})
	}
}

func TestClient_GetPostAgainstServer(t *testing.T) {
	var gotMethod, gotPath, gotAuth, gotCT, gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		gotCT = r.Header.Get("Content-Type")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c, err := NewClient(srv.URL, LogModeNone)
	require.NoError(t, err)

	t.Run("GET with option", func(t *testing.T) {
		resp, err := c.Get("/things", BearerToken("tok"))
		require.NoError(t, err)
		defer Drain(resp)
		assert.Equal(t, http.MethodGet, gotMethod)
		assert.Equal(t, "/things", gotPath)
		assert.Equal(t, "Bearer tok", gotAuth)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST JSON", func(t *testing.T) {
		resp, err := c.PostJSON("/items", map[string]any{"a": 1})
		require.NoError(t, err)
		defer Drain(resp)
		assert.Equal(t, http.MethodPost, gotMethod)
		assert.Contains(t, gotCT, "application/json")
		assert.JSONEq(t, `{"a":1}`, gotBody)
	})

	t.Run("POST form", func(t *testing.T) {
		resp, err := c.PostForm("/form", url.Values{"k": {"v"}})
		require.NoError(t, err)
		defer Drain(resp)
		assert.Equal(t, "application/x-www-form-urlencoded", gotCT)
		assert.Equal(t, "k=v", gotBody)
	})
}

func TestResponseErr(t *testing.T) {
	ok := &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(""))}
	assert.NoError(t, ResponseErr(ok))

	bad := &http.Response{
		StatusCode: 404,
		Status:     "404 Not Found",
		Body:       io.NopCloser(strings.NewReader("not found")),
	}
	err := ResponseErr(bad)
	require.Error(t, err)
	httpErr, isHTTPErr := err.(HTTPErr)
	require.True(t, isHTTPErr)
	assert.Equal(t, 404, httpErr.Code)
	assert.Equal(t, "not found", httpErr.Body)
}

func TestResponseErrOrDecodeJson(t *testing.T) {
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"name":"hexa"}`)),
	}
	var out struct {
		Name string `json:"name"`
	}
	require.NoError(t, ResponseErrOrDecodeJson(resp, &out))
	assert.Equal(t, "hexa", out.Name)

	errResp := &http.Response{
		StatusCode: 500,
		Status:     "500 Internal Server Error",
		Body:       io.NopCloser(strings.NewReader("boom")),
	}
	assert.Error(t, ResponseErrOrDecodeJson(errResp, &out))
}
