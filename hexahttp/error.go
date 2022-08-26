package hexahttp

import (
	"fmt"
	"net/http"
)

// HttpErr represents a Http response error.
type HttpErr struct {
	Code   int
	Status string
}

func (err HttpErr) Error() string {
	return fmt.Sprintf("HTTP code: %d, msg: %s, status: %s", err.Code, http.StatusText(err.Code), err.Status)
}

func responseErr(r *http.Response) error {
	if r.StatusCode <= 300 {
		return nil
	}
	return HttpErr{r.StatusCode, r.Status}
}
