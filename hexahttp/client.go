package hexahttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type ClientConfig struct {
	// optional base uri for all requests (exclude
	// absolute  http url provided in requests).
	BaseURI      *string
	AuthProvider AuthProvider
}

type RequestOption func(req *http.Request) error

// Client is improved version of the net/http client
type Client struct {
	*http.Client
	baseURI      *string
	authProvider AuthProvider
}

func New(cfg ClientConfig) *Client {
	return &Client{
		Client:       &http.Client{},
		baseURI:      cfg.BaseURI,
		authProvider: cfg.AuthProvider,
	}
}

func (c *Client) PostFormWithOptions(uri string, data url.Values, options ...RequestOption) (*http.Response, error) {
	return c.PostWithOptions(uri, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()), options...)
}

func (c *Client) PostJsonFormWithOptions(uri string, data hexa.Map, options ...RequestOption) (*http.Response, error) {
	rBody, err := json.Marshal(data)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return c.PostWithOptions(uri, "application/json;charset=UTF-8", bytes.NewBuffer(rBody), options...)
}

func (c *Client) PostWithOptions(uri string, contentType string, body io.Reader, options ...RequestOption) (*http.Response, error) {
	u, err := c.URL(uri)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	req, err := http.NewRequest("POST", u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.DoWithOptions(req, options...)
}

func (c *Client) GetWithOptions(uri string, options ...RequestOption) (*http.Response, error) {
	u, err := c.URL(uri)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return c.DoWithOptions(req, options...)
}

func (c *Client) DoWithOptions(req *http.Request, options ...RequestOption) (*http.Response, error) {
	if err := c.setOptions(req, options...); err != nil {
		return nil, tracer.Trace(err)
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return res, nil
}

func (c *Client) URL(uri string) (*url.URL, error) {
	if isValidURL(uri) {
		return url.Parse(uri)
	}

	if c.baseURI == nil {
		return nil, tracer.Trace(errors.New("client's base uri is nil, you must provide absolute url for each request"))
	}

	u, err := url.Parse(*c.baseURI)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	isAbsolutePath := len(uri) > 0 && uri[0] == '/'
	if isAbsolutePath {
		u.Path = uri
	} else {
		u.Path = path.Join(u.Path, uri)
	}

	return u, nil
}

func (c *Client) setOptions(req *http.Request, options ...RequestOption) error {
	for _, o := range options {
		if err := o(req); err != nil {
			return tracer.Trace(err)
		}
	}
	return nil
}
