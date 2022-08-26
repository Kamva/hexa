package hexahttp

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/kamva/tracer"
)

func isValidURL(val string) bool {
	u, err := url.ParseRequestURI(val)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// Bytes returns response's body bytes and close the response's body.
func Bytes(r *http.Response) ([]byte, error) {
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

// Drain drains the response body and closes its connection.
func Drain(r *http.Response) error {
	defer r.Body.Close()

	_, err := io.Copy(ioutil.Discard, r.Body)
	return err
}

// ResponseErr returns http error if the response is not successful.
// it drains the response.
func ResponseErr(r *http.Response) error {
	defer Drain(r)
	return responseErr(r)
}

// ResponseErrBytes returns the response body and its http error if
// the response isn't successful.
// It closes the response's body.
func ResponseErrBytes(r *http.Response) ([]byte, error) {
	defer r.Body.Close()
	if err := responseErr(r); err != nil {
		b, _ := io.ReadAll(r.Body)
		return b, err
	}

	_, err := io.Copy(ioutil.Discard, r.Body)
	return nil, err
}

// ResponseErrOrBytes checks if the response is successful,
// it returns the body bytes, otherwise the http error.
func ResponseErrOrBytes(r *http.Response) ([]byte, error) {
	defer r.Body.Close()

	if err := responseErr(r); err != nil {
		return nil, tracer.Trace(err)
	}

	return io.ReadAll(r.Body)
}

// ResponseErrAndBytes returns the response body byes and
// a http error if the response is not successful.
func ResponseErrAndBytes(r *http.Response) ([]byte, error) {
	defer r.Body.Close()

	b, bytesErr := Bytes(r)
	if err := responseErr(r); err != nil {
		return b, err
	}

	return b, bytesErr
}
