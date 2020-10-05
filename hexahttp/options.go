package hexahttp

import (
	"fmt"
	"github.com/kamva/hexa"
	"net/http"
)

func BasicAuth(username string, password string) RequestOption {
	return func(req *http.Request) error {
		req.SetBasicAuth(username, password)
		return nil
	}
}

func BearerToken(token string) RequestOption {
	return AuthorizationToken("Bearer", token)
}

func AuthorizationToken(tName string, token string) RequestOption {
	return func(req *http.Request) error {
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", tName, token))
		return nil
	}
}

func QueryParams(params hexa.Map) RequestOption {
	return func(req *http.Request) error {
		u := req.URL
		return UrlQueryParams(params)(u)
	}
}
