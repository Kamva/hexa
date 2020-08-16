package hexa

import (
	"errors"
	"net/http"
)

//--------------------------------
// Entity Adapter errors
//--------------------------------

var (
	ErrInvalidID = NewError(http.StatusBadRequest, "entity.invalid_id", errors.New("id value is invalid"))
)
