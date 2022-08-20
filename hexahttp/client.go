package hexahttp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	urlpkg "net/url"
	"strings"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
)

// RequestOption is like middleware which can change request before the client send it.
type RequestOption func(req *http.Request) error

// Client is improved version of the net/http client
type Client struct {
	*http.Client
	baseUrl *string
}

func NewClient(baseUrl *string) *Client {
	return &Client{
		Client:  &http.Client{},
		baseUrl: baseUrl,
	}
}

func (c *Client) PostFormWithOptions(url string, data urlpkg.Values, options ...RequestOption) (*http.Response, error) {
	return c.PostWithOptions(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()), options...)
}

func (c *Client) PostJsonFormWithOptions(urlPath string, data hexa.Map, options ...RequestOption) (*http.Response, error) {
	rBody, err := json.Marshal(data)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return c.PostWithOptions(urlPath, "application/json;charset=UTF-8", bytes.NewBuffer(rBody), options...)
}

func (c *Client) PostWithOptions(url string, contentType string, body io.Reader, options ...RequestOption) (*http.Response, error) {
	u, err := NewURL(c.baseUrl).URL(url)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	req, err := http.NewRequest("POST", u.String(), body)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	req.Header.Set("Content-Type", contentType)
	return c.DoWithOptions(req, options...)
}

func (c *Client) GetWithOptions(url string, options ...RequestOption) (*http.Response, error) {
	u, err := NewURL(c.baseUrl).URL(url)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	return c.DoWithOptions(req, options...)
}

func (c *Client) DoWithOptions(req *http.Request, options ...RequestOption) (*http.Response, error) {
	if err := c.setOptions(req, options...); err != nil {
		return nil, tracer.Trace(err)
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return res, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// TODO: improve it (do not dump if logger's debug mode is not enabled,dump response as well).
	d, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	hlog.Error("sending http request", hlog.String("request", string(d)))
	return c.Client.Do(req)
}

func (c *Client) setOptions(req *http.Request, options ...RequestOption) error {
	for _, o := range options {
		if err := o(req); err != nil {
			return tracer.Trace(err)
		}
	}
	return nil
}
