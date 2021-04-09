package hexahttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/kamva/tracer"
)

func isValidURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}

// ResponseBytes gets result of call from the client and then converts response body to bytes.
func Bytes(response *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, tracer.Trace(err)
	}

	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}

// HttpErrFromResponse return http error if response status code
// is more than 300.
func ResponseError(r *http.Response) error {
	if r.StatusCode <= 300 {
		return nil
	}

	respBytes, err := Bytes(r, nil)
	if err != nil {
		return tracer.Trace(err)
	}

	return fmt.Errorf("Http error, status: %s,  code: %d, body: %s", r.Status, r.StatusCode, string(respBytes))
}
