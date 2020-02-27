package kitty

import (
	"errors"
	"net/http"
)

//--------------------------------
// Entity Adapter errors
//--------------------------------

// Error code description:
// kt = kitty  (package or project name)
// 0 = replies about entity adapter section (identify some part in application)
// E = Error (type of code : error|response|...)
// 0 = error number zero (id of code in that part and type)

var (
	ErrInvalidID = NewError(http.StatusBadRequest, "kttmpl.2.e.0", "invalid_id", errors.New("id value is invalid"))
)
