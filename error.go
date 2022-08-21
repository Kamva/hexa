package hexa

import (
	"errors"
	"fmt"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa/hlog"
)

const (
	// ErrKeyInternalError is the internal error key in Error
	// messages over all of packages. use this to have just one
	// internal_error translation key in your translation system.
	// TODO: remove this key if we don't use it in our projects.
	ErrKeyInternalError = "lib.internal_error"
)

// Error is reply to actions when occur error in microservices.
type Error interface {
	error

	// SetError set the internal error.
	SetError(error) Error

	// InternalError returns the internal error.
	InternalError() error

	//Is function satisfy Is interface of errors package.
	Is(error) bool

	// HTTPStatus returns the http status code for the Error.
	HTTPStatus() int

	// SetHTTPStatus sets the http status code for the reply.
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
	ReportIfNeeded(hlog.Logger, Translator) bool
}

type defaultError struct {
	error

	httpStatus       int
	id               string
	localizedMessage string
	data             Map
	reportData       Map
}

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
	ee := AsHexaErr(err)
	return ee != nil && e.ID() == ee.ID()
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

func (e defaultError) ReportIfNeeded(l hlog.Logger, _ Translator) bool {
	if e.shouldReport() {
		l.With(ErrFields(e)...).Error(e.Error())
		return true
	}
	return false
}

func (e defaultError) shouldReport() bool {
	return e.HTTPStatus() >= 500
}

// NewError returns new instance the Error interface.
func NewError(httpStatus int, id string) Error {
	return defaultError{
		httpStatus: httpStatus,
		id:         id,
		//data:       make(Map), // Don't allocate memory at creation time.
		//reportData: make(Map),
	}
}

// NewLocalizedError returns new instance the Error interface.
func NewLocalizedError(status int, id string, localizedMsg string, err error) Error {
	return defaultError{
		error:            err,
		httpStatus:       status,
		id:               id,
		localizedMessage: localizedMsg,
	}
}

// AsHexaErr check whether provided error can be used as hexa error or not.
// by this method, we don't need to use guilt.CauseErr and then check if
// caused error is hexa error or not, and we can simply like any other
// error implement the errors.Unwrap interface on our Error interface.
// Why we don't use errors.As() method? We can use it also, these two
// methods do not conflict each other, it gets a new instance of target
// and check target's type and then if found error with that type in the
// chain, assign its value to the target, but here we want to just return
// our hexa error (with any implementation of hexa error) to the user.
func AsHexaErr(err error) Error {
	for err != nil {
		if hexaErr, ok := err.(Error); ok {
			return hexaErr
		}

		err = errors.Unwrap(err)
	}

	return nil
}

// Assert defaultReply implements the Error interface.
var _ Error = defaultError{}
