package kitty

import (
	"github.com/Kamva/gutil"
	"github.com/Kamva/tracer"
)

type (

	// ReplyTypeDefault is reply type of the default Reply
	ReplyTypeDefault string

	// ReplyTypeError is reply type of the Error Reply.
	ReplyTypeError string

	// ReplyParams is parameters of the reply to use in translation,...
	ReplyParams map[string]interface{}

	// ReplyData is extra data of the reply to show to the user.
	ReplyData map[string]interface{}

	// ReplyReportData is the report data that use on reporting reply to somewhere (e.g log aggregator).
	ReplyReportData map[string]interface{}

	// Reply is reply to actions in microservices.
	// Important : dont forget to implement setter functions
	// in all of Reply implmentaions, not just base defaultReply
	// struct, because on return Reply instace from setter methods
	// defaultReply struct return instance of itself, so you sould
	// implement yourself setter functions for your struct to prevent
	// converting from your struct instance to defaultReply instance in
	// setter methods that implemented by defaultReply.
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

		// Type returns the reply type, you can use type assertion
		// to detect the reply type.
		// e.g _,ok:=r.Type().(kitty.ReplyTypeDefault)
		Type() interface{}

		// ReportIfNeeded function report the reply to the log system if
		// http status code is in range 5XX.
		// return value specify that reported or no.
		ReportIfNeeded(Logger, Translator) bool

		// HTTPStatus returns the http status code for the reply.
		HTTPStatus() int

		// HTTPStatus returns the http status code for the reply.
		SetHTTPStatus(status int) Reply

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

		// Data returns the extra data of the reply (e.g show this data to user).
		Data() ReplyData

		// SetData set the reply data as extra data of the reply to show to the user.
		SetData(data ReplyData) Reply

		// ReportData returns the data that should use on reporting reply to somewhere (e.g log aggregator)
		ReportData() ReplyReportData

		SetReportData(data ReplyReportData) Reply
	}

	// Alias Reply as Error
	Error Reply

	// defaultReply implements the Reply interface.
	defaultReply struct {
		replyType   interface{}
		httpStatus  int
		code        string
		key         string
		internalMsg string
		params      ReplyParams
		data        ReplyData
		reportData  ReplyReportData
	}

	defaultError struct {
		defaultReply
	}
)

const (
	// ErrorKeyInternalError is the internal error key in reply
	// messages over all of packages. use this to have just one
	// internal_error translation key in your translation system.
	ReplyErrKeyInternalError = "internal_error"
)

func (r defaultReply) Is(err error) bool {
	e, ok := tracer.Cause(err).(Reply)
	return ok && r.Code() == e.Code()
}

func (r defaultReply) Type() interface{} {
	return r.replyType
}

// Error implements to just satisfy the Reply interface.
func (r defaultReply) Error() string {
	return "__default_reply__"
}

func (r defaultReply) InternalMessage() string {
	return r.internalMsg
}

func (r defaultReply) SetInternalMessage(msg string) Reply {
	r.internalMsg = msg
	return r
}

func (r defaultReply) SetInternalMsg() string {
	return r.internalMsg
}

func (r defaultReply) shouldReport() bool {
	return r.HTTPStatus() >= 500
}

func (r defaultReply) ReportIfNeeded(l Logger, t Translator) bool {
	if r.shouldReport() {
		l.WithFields(reportFields(r)...).Info(r.Error())
		return true
	}

	return false
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

func (r defaultReply) Params() ReplyParams {
	return r.params
}

func (r defaultReply) SetParams(params ReplyParams) Reply {
	r.params = params
	return r
}

func (r defaultReply) Data() ReplyData {
	return r.data
}

func (r defaultReply) SetData(data ReplyData) Reply {
	r.data = data
	return r
}

func (r defaultReply) ReportData() ReplyReportData {
	return r.reportData
}

func (r defaultReply) SetReportData(data ReplyReportData) Reply {
	r.reportData = data
	return r
}

// Error method returns the error message.
func (e defaultError) Error() string {
	return e.internalMsg
}

func (e defaultError) ReportIfNeeded(l Logger, t Translator) bool {
	if e.shouldReport() {
		l.WithFields(reportFields(e)...).Error(e.Error())
		return true
	}

	return false
}

func (e defaultError) SetHTTPStatus(status int) Reply {
	e.httpStatus = status

	return e
}

func (e defaultError) SetInternalMessage(msg string) Reply {
	e.internalMsg = msg
	return e
}

func (e defaultError) SetParams(params ReplyParams) Reply {
	e.params = params
	return e
}

func (e defaultError) SetData(data ReplyData) Reply {
	e.data = data
	return e
}

func (e defaultError) SetReportData(data ReplyReportData) Reply {
	e.reportData = data
	return e
}

// NewReply returns new instance the Reply interface implemented by defaultReply.
func NewReply(httpStatus int, code string, key string, iMsg string) Reply {
	return defaultReply{
		replyType:   ReplyTypeDefault("__reply_default__"),
		httpStatus:  httpStatus,
		code:        code,
		key:         key,
		internalMsg: iMsg,
		params:      make(ReplyParams),
		data:        make(ReplyData),
		reportData:  make(ReplyReportData),
	}
}

// NewReply returns new instance the Reply interface implemented by defaultReply.
func NewError(httpStatus int, code string, key string, err string) Error {
	return defaultError{
		defaultReply{
			replyType:   ReplyTypeError("__reply_error__"),
			httpStatus:  httpStatus,
			code:        code,
			key:         key,
			internalMsg: err,
			params:      make(ReplyParams),
			data:        make(ReplyData),
			reportData:  make(ReplyReportData),
		},
	}
}

// StructToReplyData convert struct to reply data
func StructToReplyData(input interface{}) ReplyData {
	return gutil.StructToMap(input)
}

// reportFields return fields that need to include in reply report.
func reportFields(r Reply) []interface{} {
	data := map[string]interface{}{
		"__type__":        r.Type(),
		"__code__":        r.Code(),
		"__http_status__": r.HTTPStatus(),
		"__data__":        r.Data(),
	}

	fields := append(gutil.MapToKeyValue(data), gutil.MapToKeyValue(r.ReportData())...)

	return fields
}

// Assert defaultReply implements the Error interface.
var _ Reply = defaultReply{}
