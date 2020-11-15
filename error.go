package hexa

import (
	"fmt"
	"github.com/kamva/gutil"
	"github.com/kamva/tracer"
)

type (
	// Error is reply to actions when occur error in microservices.
	Error interface {
		error

		// SetError set the internal error.
		SetError(error) Error

		// InternalError returns the internal error.
		InternalError() error

		//Is function satisfy Is interface of errors package.
		Is(error) bool

		// HTTPStatus returns the http status code for the Error.
		HTTPStatus() int

		// HTTPStatus returns the http status code for the reply.
		SetHTTPStatus(status int) Error

		// ID is error's identifier. its format should be
		// something like "product.variant.not_found" or "lib.jwt.not_found" .
		// as a convention we prefix our base packages (all hexa packages) with "lib".
		ID() string

		// Localize localize te message for you.
		// you can store the gRPC localized error
		// message and return it by this method.
		Localize(t Translator) (string, error)

		// Data returns the extra data of the Error (e.g show this data to user).
		// Note: we use data as translation prams also.
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

		httpStatus       int
		id               string
		localizedMessage string
		data             Map
		reportData       Map
	}
)

const (
	// ErrorKeyInternalError is the internal error key in Error
	// messages over all of packages. use this to have just one
	// internal_error translation key in your translation system.
	// TODO: remove this key if we don't use it in our projects.
	ErrKeyInternalError = "lib.internal_error"
)

func (e defaultError) Error() string {
	if e.error != nil {
		return e.error.Error()
	}

	return fmt.Sprintf("Error with id: %s", e.ID())
}

func (e defaultError) SetError(err error) Error {
	e.error = err
	return e
}

func (e defaultError) InternalError() error {
	return e.error
}

func (e defaultError) Is(err error) bool {
	ee, ok := gutil.CauseErr(err).(Error)
	return ok && e.ID() == ee.ID()
}

func (e defaultError) HTTPStatus() int {
	return e.httpStatus
}

func (e defaultError) SetHTTPStatus(status int) Error {
	e.httpStatus = status
	return e
}

func (e defaultError) ID() string {
	return e.id
}

func (e defaultError) Localize(t Translator) (string, error) {
	if e.localizedMessage != "" {
		return e.localizedMessage, nil
	}
	return t.Translate(e.ID(), gutil.MapToKeyValue(e.Data())...)
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
			"__error_id__":    e.ID(),
			"__http_status__": e.HTTPStatus(),
			"__data__":        e.Data(),
			"__report__":      e.ReportData(),
		}

		// If exists error and error is traced,print its stack.
		if stack := tracer.StackAsString(tracer.MoveStackIfNeeded(e, e.error)); stack != "" {
			data[ErrorStackLogKey] = stack
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
func NewError(httpStatus int, id string, err error) Error {
	return defaultError{
		error:      err,
		httpStatus: httpStatus,
		id:         id,
		data:       make(Map),
		reportData: make(Map),
	}
}

// NewError returns new instance the Error interface.
func NewLocalizedError(status int, id string, localizedMsg string, err error) Error {
	return defaultError{
		error:            err,
		httpStatus:       status,
		id:               id,
		localizedMessage: localizedMsg,
		data:             make(Map),
		reportData:       make(Map),
	}
}

// Assert defaultReply implements the Error interface.
var _ Error = defaultError{}
