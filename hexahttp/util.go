package hexahttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/kamva/tracer"
)

func isValidURL(val string) bool {
	u, err := url.ParseRequestURI(val)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// Bytes gets result of call from the client and then converts response body to bytes.
func Bytes(response *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, tracer.Trace(err)
	}

	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}

// ResponseErr returns http error if response status code
// is more than 300.
func ResponseErr(r *http.Response, err error) (*http.Response, error) {
	if err != nil {
		return r, tracer.Trace(err)
	}

	if r.StatusCode <= 300 {
		return r, nil
	}

	return r, fmt.Errorf("Http response error, status: %s,  code: %d", r.Status, r.StatusCode)
}

// ResponseErrOrBytes checks if the response doesn't have any error and error and
// its HTTP status code is a 2xx response status code, returns the response and its
// body content bytes.
func ResponseErrOrBytes(r *http.Response, err error) (*http.Response, []byte, error) {
	if r, err = ResponseErr(r, err); err != nil {
		return r, nil, tracer.Trace(err)
	}

	b, err := Bytes(r, err)
	if err != nil {
		return r, nil, tracer.Trace(err)
	}

	return r, b, nil
}
