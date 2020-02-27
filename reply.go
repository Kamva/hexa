package kitty

import (
	"github.com/Kamva/tracer"
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

		// Code return the reply identifier code
		Code() string

		// Key returns unique key for each reply to use as translation key,...
		Key() string

		// Params returns params of the reply to use in translation,...
		Params() Map

		// SetParams set the reply parameters to use in reply translation,...
		SetParams(params Map) Reply

		// Data returns the extra data of the reply (e.g show this data to user).
		Data() Map

		// SetData set the reply data as extra data of the reply to show to the user.
		SetData(data Map) Reply
	}

	// defaultReply implements the Reply interface.
	defaultReply struct {
		httpStatus int
		code       string
		key        string
		params     Map
		data       Map
	}
)

func (r defaultReply) Is(err error) bool {
	e, ok := tracer.Cause(err).(Reply)
	return ok && r.Code() == e.Code()
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

func (r defaultReply) Code() string {
	return r.code
}

func (r defaultReply) Key() string {
	return r.key
}

func (r defaultReply) Params() Map {
	return r.params
}

func (r defaultReply) SetParams(params Map) Reply {
	r.params = params
	return r
}

func (r defaultReply) Data() Map {
	return r.data
}

func (r defaultReply) SetData(data Map) Reply {
	r.data = data
	return r
}

// NewReply returns new instance the Reply interface implemented by defaultReply.
func NewReply(httpStatus int, code string, key string) Reply {
	return defaultReply{
		httpStatus: httpStatus,
		code:       code,
		key:        key,
		params:     make(Map),
		data:       make(Map),
	}
}


// Assert defaultReply implements the Error interface.
var _ Reply = defaultReply{}
