package hexa

type (
	// Reply is reply to actions in microservices.
	Reply interface {
		// HTTPStatus returns the http status code for the reply.
		HTTPStatus() int

		// SetHTTPStatus sets the http status code for the reply.
		SetHTTPStatus(status int) Reply

		// ID is reply identifier
		ID() string

		// Data returns the extra data of the reply (e.g show this data to user).
		// Note: we use data as translation prams also.
		Data() any

		// SetData set the reply data as extra data of the reply to show to the user.
		SetData(data any) Reply
	}

	// defaultReply implements the Reply interface.
	defaultReply struct {
		httpStatus int
		id         string
		data       any
	}
)

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

func (r defaultReply) Data() any {
	return r.data
}

func (r defaultReply) SetData(data any) Reply {
	r.data = data
	return r
}

// NewReply returns new instance the Reply interface implemented by defaultReply.
func NewReply(httpStatus int, id string) Reply {
	return defaultReply{
		httpStatus: httpStatus,
		id:         id,
	}
}

// Assert defaultReply implements the Error interface.
var _ Reply = defaultReply{}
