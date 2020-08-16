package hexa

import (
	"github.com/kamva/gutil"
)

type (
	// Reply is reply to actions in microservices.
	Reply interface {
		// specifically this interface contains the error interface
		// also, to able pass as return value in some frameworks that
		// permit to return error and then set error handler.
		// e.g by this  way we can pass this Reply to the ErrorHandler
		// and error handler can check Reply type and return the proper
		// response.
		// in implementation of Reply we can return just empty string or
		// special value (like "__reply__") as error message.
		error

		//Is function satisfy Is interface of errors package.
		Is(error) bool

		// HTTPStatus returns the http status code for the reply.
		HTTPStatus() int

		// HTTPStatus returns the http status code for the reply.
		SetHTTPStatus(status int) Reply

		// ID is reply identifier
		ID() string

		// Data returns the extra data of the reply (e.g show this data to user).
		// Note: we use data as translation prams also.
		Data() Map

		// SetData set the reply data as extra data of the reply to show to the user.
		SetData(data Map) Reply
	}

	// defaultReply implements the Reply interface.
	defaultReply struct {
		httpStatus int
		id         string
		data       Map
	}
)

func (r defaultReply) Is(err error) bool {
	e, ok := gutil.CauseErr(err).(Reply)
	return ok && r.ID() == e.ID()
}

// Error implements to just satisfy the Reply and error interface.
func (r defaultReply) Error() string {
	return "__reply__"
}

func (r defaultReply) HTTPStatus() int {
	return r.httpStatus
}

func (r defaultReply) SetHTTPStatus(status int) Reply {
	r.httpStatus = status

	return r
}

func (r defaultReply) ID() string {
	return r.id
}

func (r defaultReply) Data() Map {
	return r.data
}

func (r defaultReply) SetData(data Map) Reply {
	r.data = data
	return r
}

// NewReply returns new instance the Reply interface implemented by defaultReply.
func NewReply(httpStatus int, id string) Reply {
	return defaultReply{
		httpStatus: httpStatus,
		id:         id,
		data:       make(Map),
	}
}

// Assert defaultReply implements the Error interface.
var _ Reply = defaultReply{}
