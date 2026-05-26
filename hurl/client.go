package hurl

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	urlpkg "net/url"
	"strings"

	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
)

const (
	LogModeNone uint = 0

	LogModeRequestHeaders uint = 1 << iota
	LogModeRequestBody
	LogModeResponseHeaders
	LogModeResponseBody

	LogModeRequest  = LogModeRequestHeaders | LogModeRequestBody
	LogModeResponse = LogModeResponseHeaders | LogModeResponseBody
	LogModeAll      = LogModeRequest | LogModeResponse
)

// RequestOption is like middleware which can change request before the client send it.
type RequestOption func(req *http.Request) error

// Client is improved version of the net/http client
type Client struct {
	*http.Client
	base    *URL
	logMode uint
}

func NewClient(base string, logMode uint) (*Client, error) {
	// NewURL returns a non-nil *URL even for an empty base, so methods on
	// c.base never dereference a nil receiver. A relative request path with
	// no base then returns a clear error instead of panicking.
	u, err := NewURL(base)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return NewClientWithOptions(&http.Client{}, u, logMode), nil
}

// NewClientWithOptions creates a new HTTP client. base url
// could be nil. in that case all request urls
// must be absolute.
func NewClientWithOptions(cli *http.Client, base *URL, logMode uint) *Client {
	return &Client{
		Client:  cli,
		base:    base,
		logMode: logMode,
	}
}

func (c *Client) Head(url string, options ...RequestOption) (resp *http.Response, err error) {
	u, err := c.base.URL(url)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	req, err := http.NewRequest("HEAD", u.String(), nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req, options...)
}

func (c *Client) PostForm(url string, data urlpkg.Values, options ...RequestOption) (*http.Response, error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()), options...)
}

func (c *Client) PostJSON(urlPath string, data any, options ...RequestOption) (*http.Response, error) {
	rBody, err := json.Marshal(data)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return c.Post(urlPath, "application/json;charset=UTF-8", bytes.NewBuffer(rBody), options...)
}

func (c *Client) Post(url string, contentType string, body io.Reader, options ...RequestOption) (*http.Response, error) {
	u, err := c.base.URL(url)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	req, err := http.NewRequest("POST", u.String(), body)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req, options...)
}

func (c *Client) Get(url string, options ...RequestOption) (*http.Response, error) {
	u, err := c.base.URL(url)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	return c.Do(req, options...)
}

func (c *Client) Do(req *http.Request, options ...RequestOption) (*http.Response, error) {
	// Apply options
	for _, o := range options {
		if err := o(req); err != nil {
			return nil, tracer.Trace(err)
		}
	}

	if c.logMode&LogModeRequestHeaders != 0 {
		// Redact sensitive headers for the dump only, then restore the real
		// ones before sending so credentials aren't leaked to logs.
		origHeader := req.Header
		req.Header = redactHeaders(origHeader)
		dumped, err := httputil.DumpRequestOut(req, c.logMode&LogModeRequestBody != 0)
		req.Header = origHeader
		if err != nil {
			return nil, tracer.Trace(err)
		}

		hlog.Info("outgoing HTTP request", hlog.String("request", string(dumped)))
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	if c.logMode&LogModeResponseHeaders != 0 {
		origHeader := resp.Header
		resp.Header = redactHeaders(origHeader)
		dumped, err := httputil.DumpResponse(resp, c.logMode&LogModeResponseBody != 0)
		resp.Header = origHeader
		if err != nil {
			return nil, tracer.Trace(err)
		}

		hlog.Info("incoming HTTP response", hlog.String("response", string(dumped)))
	}

	return resp, nil
}

// sensitiveHeaders are request/response headers whose values are masked
// before a message is dumped to logs, so credentials and session data don't
// leak when LogMode*Headers is enabled.
var sensitiveHeaders = []string{
	"Authorization",
	"Proxy-Authorization",
	"Cookie",
	"Set-Cookie",
}

// redactHeaders returns a copy of h with the values of sensitiveHeaders
// replaced by "REDACTED". The input header is left unmodified. It is nil-safe.
func redactHeaders(h http.Header) http.Header {
	if h == nil {
		return nil
	}

	clone := h.Clone()
	for _, name := range sensitiveHeaders {
		if len(clone.Values(name)) != 0 {
			clone.Set(name, "REDACTED")
		}
	}
	return clone
}
