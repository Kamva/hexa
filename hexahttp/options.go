package hexahttp

import (
	"fmt"
	"net/http"
	urlpkg "net/url"

	"github.com/kamva/hexa"
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

func AuthorizationToken(tokenType string, token string) RequestOption {
	return AuthenticateHeader("Authorization", tokenType, token)
}

func AuthenticateHeader(header string, tokenType string, token string) RequestOption {
	return func(req *http.Request) error {
		val := fmt.Sprintf("%v %v", tokenType, token)
		if tokenType == "" {
			val = token
		}

		req.Header.Set(header, val)
		return nil
	}
}

func QueryParams(params hexa.Map) RequestOption {
	return func(req *http.Request) error {
		u := req.URL
		return URLQueryParams(params)(u)
	}
}

//--------------------------------
// URL options
//--------------------------------

func URLQueryParams(params hexa.Map) URLOption {
	return func(u *urlpkg.URL) error {
		q := u.Query()

		for k, v := range params {
			q.Add(k, fmt.Sprint(v))
		}
		u.RawQuery = q.Encode()
		return nil
	}
}
