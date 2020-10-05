package hexahttp

import (
	"errors"
	"fmt"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	urlpkg "net/url"
	"path"
)

type URLOption func(url *urlpkg.URL) error

type URL struct {
	baseUrl *string
}

func NewURL(baseUrl *string) *URL {
	return &URL{
		baseUrl: baseUrl,
	}
}

func (u *URL) URL(url string, options ...URLOption) (resultUrl *urlpkg.URL, err error) {
	if isValidURL(url) {
		resultUrl, err = urlpkg.Parse(url)
	} else {
		resultUrl, err = u.getFromBase(url)
	}
	if err != nil {
		return nil, tracer.Trace(err)
	}

	for _, o := range options {
		if err := o(resultUrl); err != nil {
			return nil, tracer.Trace(err)
		}
	}

	return resultUrl, nil
}

func (u *URL) getFromBase(urlPath string) (*urlpkg.URL, error) {
	if u.baseUrl == nil {
		return nil, tracer.Trace(errors.New("provide base url otherwise you need to send the full url"))
	}

	resultUrl, err := urlpkg.Parse(*u.baseUrl)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	isAbsolutePath := len(urlPath) > 0 && urlPath[0] == '/'
	if isAbsolutePath {
		resultUrl.Path = urlPath
	} else {
		resultUrl.Path = path.Join(resultUrl.Path, urlPath)
	}

	return resultUrl, nil
}

func UrlQueryParams(params hexa.Map) URLOption {
	return func(u *urlpkg.URL) error {
		q := u.Query()

		for k, v := range params {
			q.Add(k, fmt.Sprint(v))
		}
		u.RawQuery = q.Encode()
		return nil
	}
}
