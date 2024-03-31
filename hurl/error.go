package hurl

import (
	"fmt"
	"net/http"
)

// HTTPErr represents a Http response error.
type HTTPErr struct {
	Code   int
	Status string
	Body string
}

func (err HTTPErr) Error() string {
	return fmt.Sprintf("HTTP code: %d, msg: %s, status: %s, body: %s", err.Code, http.StatusText(err.Code), err.Status, err.Body)
}

// ResponseErr returns http error if the response is not successful.
func responseErr(r *http.Response) error {
	if r.StatusCode <= 300 {
		return nil
	}

	body, _ := io.ReadAll(r.Body)

	return HttpErr{r.StatusCode, r.Status, string(body)}
}
