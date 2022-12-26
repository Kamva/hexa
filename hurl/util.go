package hurl

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

func isValidURL(val string) bool {
	u, err := url.ParseRequestURI(val)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// Bytes returns response's body bytes and close the response's body.
//func Bytes(r *http.Response) ([]byte, error) {
//	defer r.Body.Close()
//	return io.ReadAll(r.Body)
//}

// Drain drains the response body and closes its connection.
func Drain(r *http.Response) error {
	if r.Body != nil {
		_, _ = io.Copy(io.Discard, r.Body)
		return r.Body.Close()
	}
	return nil
}

// ResponseErrBytes returns the response body and its http error if
// the response isn't successful.
// It closes the response's body.
//func ResponseErrBytes(r *http.Response) ([]byte, error) {
//	defer Drain(r) //nolint
//	if err := ResponseErr(r); err != nil {
//		b, _ := io.ReadAll(r.Body)
//		return b, err
//	}
//	return nil, nil
//}

// ResponseErrOrBytes checks if the response is successful,
// it returns the body bytes, otherwise the http error.
//func ResponseErrOrBytes(r *http.Response) ([]byte, error) {
//	defer Drain(r) //nolint
//
//	if err := ResponseErr(r); err != nil {
//		return nil, tracer.Trace(err)
//	}
//
//	return io.ReadAll(r.Body)
//}

func ResponseErrOrDecodeJson(r *http.Response, val interface{}) error { //nolint:revive
	if err := ResponseErr(r); err != nil {
		return err
	}

	return json.NewDecoder(r.Body).Decode(val)
}

// ResponseErrAndBytes returns the response body bytes and
// a http error if the response is not successful.
//func ResponseErrAndBytes(r *http.Response) ([]byte, error) {
//	b, bytesErr := Bytes(r)
//	if err := ResponseErr(r); err != nil {
//		return b, err
//	}
//
//	return b, bytesErr
//}
