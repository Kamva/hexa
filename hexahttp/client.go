package hexahttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"io"
	"net/http"
	urlpkg "net/url"
	"path"
	"strings"
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
	u, err := c.URL(url)
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

func (c *Client) GetWithOptions(url string, options ...RequestOption) (*http.Response, error) {
	u, err := c.URL(url)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return c.DoWithOptions(req, options...)
}

func (c *Client) DoWithOptions(url *http.Request, options ...RequestOption) (*http.Response, error) {
	if err := c.setOptions(url, options...); err != nil {
		return nil, tracer.Trace(err)
	}

	res, err := c.Client.Do(url)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return res, nil
}

func (c *Client) URL(url string) (*urlpkg.URL, error) {
	if isValidURL(url) {
		return urlpkg.Parse(url)
	}

	if c.baseUrl == nil {
		return nil, tracer.Trace(errors.New("provide client base url otherwise client needs absolute url for each request"))
	}

	u, err := urlpkg.Parse(*c.baseUrl)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	isAbsolutePath := len(url) > 0 && url[0] == '/'
	if isAbsolutePath {
		u.Path = url
	} else {
		u.Path = path.Join(u.Path, url)
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
