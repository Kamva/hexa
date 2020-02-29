package kitty

import (
	"fmt"
	"github.com/Kamva/gutil"
	"github.com/Kamva/tracer"
)

type (
	// Error is reply to actions when occur error in microservices.
	Error interface {
		error

		// SetError set the internal error.
		SetError(error) Error

		//Is function satisfy Is interface of errors package.
		Is(error) bool

		// HTTPStatus returns the http status code for the Error.
		HTTPStatus() int

		// HTTPStatus returns the http status code for the reply.
		SetHTTPStatus(status int) Error

		// Code return the Error identifier code
		Code() string

		// Key returns unique key for each Error to use as translation key,...
		Key() string

		// Params returns params of the Error to use in translation,...
		Params() Map

		// SetParams set the Error translation parameters to use in reply translation,...
		SetParams(params Map) Error

		// Data returns the extra data of the Error (e.g show this data to user).
		Data() Map

		// SetData set the Error data as extra data of the Error to show to the user.
		SetData(data Map) Error

		// ReportData returns the data that should use on reporting Error to somewhere (e.g log aggregator)
		ReportData() Map

		SetReportData(data Map) Error

		// ReportIfNeeded function report the Error to the log system if
		// http status code is in range 5XX.
		// return value specify that reported or no.
		ReportIfNeeded(Logger, Translator) bool
	}

	defaultError struct {
		error

		httpStatus int
		code       string
		key        string
		params     Map
		data       Map
		reportData Map
	}
)

const (
	// ErrorKeyInternalError is the internal error key in Error
	// messages over all of packages. use this to have just one
	// internal_error translation key in your translation system.
	ErrKeyInternalError = "err_internal_error"
)

func (e defaultError) Error() string {
	if e.error != nil {
		return e.error.Error()
	}

	return ""
}

func (e defaultError) SetError(err error) Error {
	e.error = err
	return e
}

func (e defaultError) Is(err error) bool {
	ee, ok := tracer.Cause(err).(Error)
	return ok && e.Code() == ee.Code()
}

func (e defaultError) HTTPStatus() int {
	return e.httpStatus
}

func (e defaultError) SetHTTPStatus(status int) Error {
	e.httpStatus = status
	return e
}

func (e defaultError) Code() string {
	return e.code
}

func (e defaultError) Key() string {
	return e.key
}

func (e defaultError) Params() Map {
	return e.params
}

func (e defaultError) SetParams(params Map) Error {
	e.params = params
	return e
}

func (e defaultError) Data() Map {
	return e.data
}

func (e defaultError) SetData(data Map) Error {
	e.data = data
	return e
}

func (e defaultError) ReportData() Map {
	return e.reportData
}

func (e defaultError) SetReportData(data Map) Error {
	e.reportData = data
	return e
}

func (e defaultError) ReportIfNeeded(l Logger, t Translator) bool {
	if e.shouldReport() {
		data := map[string]interface{}{
			"__code__":        e.Code(),
			"__http_status__": e.HTTPStatus(),
			"__data__":        e.Data(),
			"__report__":      e.ReportData(),
		}

		// If exists error and error is traced,print its stack.
		if te, ok := e.error.(tracer.StackTracer); e.error != nil && ok {
			data["__trace__"] = fmt.Sprintf("%+v", te)
		}

		fields := append(gutil.MapToKeyValue(data), gutil.MapToKeyValue(e.ReportData())...)
		l.WithFields(fields...).Error(e.Error())
		return true
	}
	return false
}

func (e defaultError) shouldReport() bool {
	return e.HTTPStatus() >= 500
}

// NewError returns new instance the Error interface.
func NewError(httpStatus int, code string, key string, err error) Error {
	return defaultError{
		error:      err,
		httpStatus: httpStatus,
		code:       code,
		key:        key,
		params:     make(Map),
		data:       make(Map),
		reportData: make(Map),
	}
}

// Assert defaultReply implements the Error interface.
var _ Error = defaultError{}
