package kitty

import "github.com/Kamva/gutil"

type (

	// ReplyTypeDefault is reply type of the default Reply
	ReplyTypeDefault string

	// ReplyTypeError is reply type of the Error Reply.
	ReplyTypeError string

	// ReplyData is extra data of the reply to show to the user.
	ReplyData map[string]interface{}

	// ReplyParams is parameters of the reply to use in translation,...
	ReplyParams map[string]interface{}

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

		// Type returns the reply type, you can use type assertion
		// to detect the reply type.
		// e.g _,ok:=r.Type().(kitty.ReplyTypeDefault)
		Type() interface{}

		// ShouldReport method specify that whether reply should report to the log system or no.
		ShouldReport() bool

		// Report function report the reply to the log system.
		Report(Logger, Translator)

		// HTTPStatus returns the http status code for the reply.
		HTTPStatus() int

		// InternalMessage returns the internal message.
		InternalMessage() string

		// SetInternalMessage set the internal message (e.g to report to log system)
		SetInternalMessage(msg string) Reply

		// Code return the reply identifier code
		Code() string

		// Key returns unique key for each reply to use as translation key,...
		Key() string

		// Params returns params of the reply to use in translation,...
		Params() ReplyParams

		// SetParams set the reply parameters to use in reply translation,...
		SetParams(params ReplyParams) Reply

		// Data returns the extra data of the reply.
		Data() ReplyData

		// SetData set the reply data as extra data of the reply to show to the user.
		SetData(data ReplyData) Reply
	}

	// Alias Reply as Error
	Error Reply

	// defaultReply implements the Reply interface.
	defaultReply struct {
		replyType    interface{}
		shouldReport bool
		httpStatus   int
		code         string
		key          string
		internalMsg  string
		params       ReplyParams
		data         ReplyData
	}

	defaultError struct {
		defaultReply
	}
)

const (
	// ErrorKeyInternalError is the internal error key in reply
	//messages over all of packages. use this to have just one
	// internal_error translation key in your translation system.
	ErrorKeyInternalError = "internal_error"
)

func (e defaultReply) Type() interface{} {
	return e.replyType
}

// Error implements to just satisfy the Reply interface.
func (e defaultReply) Error() string {
	return "__default_reply___"
}

func (e defaultReply) InternalMessage() string {
	return e.internalMsg
}
func (e defaultReply) SetInternalMessage(msg string) Reply {
	e.internalMsg = msg
	return e
}

func (e defaultReply) SetInternalMsg() string {
	return e.internalMsg
}

func (e defaultReply) ShouldReport() bool {
	return e.shouldReport
}

func (e defaultReply) Report(l Logger, t Translator) {
	data := map[string]interface{}{
		"__type__":        e.replyType,
		"__code__":        e.Code(),
		"__http_status__": e.HTTPStatus(),
	}

	fields := append(gutil.MapToKeyValue(data), gutil.MapToKeyValue(e.Data())...)

	l.WithFields(fields...).Info(e.Error())
}

func (e defaultReply) HTTPStatus() int {
	return e.httpStatus
}

func (e defaultReply) Code() string {
	return e.code
}

func (e defaultReply) Key() string {
	return e.code
}

func (e defaultReply) Params() ReplyParams {
	return e.params
}

func (e defaultReply) SetParams(params ReplyParams) Reply {
	e.params = params
	return e
}

func (e defaultReply) Data() ReplyData {
	return e.data
}

func (e defaultReply) SetData(data ReplyData) Reply {
	e.data = data
	return e
}

// Error method returns the error message.
func (e defaultError) Error() string {
	return e.internalMsg
}

func (e defaultError) Report(l Logger, t Translator) {
	data := map[string]interface{}{
		"__type__":        e.replyType,
		"__code__":        e.Code(),
		"__http_status__": e.HTTPStatus(),
	}

	fields := append(gutil.MapToKeyValue(data), gutil.MapToKeyValue(e.Data())...)

	l.WithFields(fields...).Error(e.Error())
}

// NewReply returns new instance the Reply interface implemented by defaultReply.
func NewReply(shouldReport bool, httpStatus int, code string, key string, iMsg string) Reply {
	return defaultReply{
		replyType:    ReplyTypeDefault("__reply_default__"),
		shouldReport: shouldReport,
		httpStatus:   httpStatus,
		code:         code,
		key:          key,
		internalMsg:  iMsg,
	}
}

// NewReply returns new instance the Reply interface implemented by defaultReply.
func NewError(shouldReport bool, httpStatus int, code string, key string, err string) Error {
	return defaultError{
		defaultReply{
			replyType:    ReplyTypeDefault("__reply_error__"),
			shouldReport: shouldReport,
			httpStatus:   httpStatus,
			code:         code,
			key:          key,
			internalMsg:  err,
		},
	}
}

// Assert defaultReply implements the Error interface.
var _ Reply = defaultReply{}
