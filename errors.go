package hexa

import (
	"errors"
	"net/http"
)

//--------------------------------
// Entity Adapter errors
//--------------------------------

var (
	ErrInvalidID = NewError(http.StatusBadRequest, "lib.entity.invalid_id").SetError(errors.New("id value is invalid"))
)
