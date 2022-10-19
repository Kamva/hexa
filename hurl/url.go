package hurl

import (
	"errors"
	urlpkg "net/url"
	"path"

	"github.com/kamva/tracer"
)

type URLOption func(url *urlpkg.URL) error

type URL struct {
	base *urlpkg.URL
}

// NewURL creates a new url with the base url value.
// we can use it to convert relative url paths to
// absolute path.
func NewURL(base string) (*URL, error) {
	if base == "" {
		return new(URL), nil
	}

	parsed, err := urlpkg.Parse(base)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return &URL{
		base: parsed,
	}, nil
}

func (u *URL) URL(url string, options ...URLOption) (resultURL *urlpkg.URL, err error) {
	if isValidURL(url) {
		resultURL, err = urlpkg.Parse(url)
	} else {
		resultURL, err = u.getFromBase(url)
	}
	if err != nil {
		return nil, tracer.Trace(err)
	}

	for _, o := range options {
		if err := o(resultURL); err != nil {
			return nil, tracer.Trace(err)
		}
	}

	return resultURL, nil
}

func (u *URL) getFromBase(urlPath string) (*urlpkg.URL, error) {
	if u.base == nil {
		return nil, tracer.Trace(errors.New("provide base url otherwise you need to send the full url"))
	}

	resultURL := u.copyBase()

	isAbsolutePath := len(urlPath) > 0 && urlPath[0] == '/'
	if isAbsolutePath {
		resultURL.Path = urlPath
	} else {
		resultURL.Path = path.Join(resultURL.Path, urlPath)
	}

	return resultURL, nil
}

func (u *URL) copyBase() *urlpkg.URL {
	if u.base == nil {
		return nil
	}
	cpy := *u.base
	return &cpy
}
