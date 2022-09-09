package hexahttp

import (
	"fmt"
	"net/http"
)

// HTTPErr represents a Http response error.
type HTTPErr struct {
	Code   int
	Status string
}

func (err HTTPErr) Error() string {
	return fmt.Sprintf("HTTP code: %d, msg: %s, status: %s", err.Code, http.StatusText(err.Code), err.Status)
}

func responseErr(r *http.Response) error {
	if r.StatusCode <= 300 {
		return nil
	}
	return HTTPErr{r.StatusCode, r.Status}
}
