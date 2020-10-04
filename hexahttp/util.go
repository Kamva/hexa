package hexahttp

import (
	"github.com/kamva/tracer"
	"io/ioutil"
	"net/http"
	"net/url"
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
