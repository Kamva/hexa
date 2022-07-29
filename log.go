package hexa

import (
	"context"

	"github.com/kamva/tracer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const ErrorStackLogKey = "_stack"

// LogFunc gets logger, set data on it and return it.
// example of such functions is ErrStack function in the
// hlog package.
type LogFunc func(Logger) Logger

// Type and function aliases from zap to limit the libraries scope into hexa code

type LogField = zapcore.Field

// LogField helpers to create fields from types.

var Int64Field = zap.Int64
var Int32Field = zap.Int32
var IntField = zap.Int
var Uint32Field = zap.Uint32
var Uint64Field = zap.Uint64
var StringField = zap.String
var AnyField = zap.Any
var ErrField = zap.Error
var NamedErrField = zap.NamedError
var BoolField = zap.Bool
var DurationField = zap.Duration
var TimeField = zap.Time
var TimesField = zap.Times
var TimepField = zap.Timep

// ErrStackField print error stack(if exists) using logger
func ErrStackField(err error) LogField {
	return StringField(ErrorStackLogKey, tracer.StackAsString(err))
}

// ErrFields checks if the provided error is a Hexa error, returns
// hexa error fields, otherwise returns regular error fields.
func ErrFields(err error) []LogField {
	if hexaErrFields := hexaErrFields(err); len(hexaErrFields) != 0 {
		return hexaErrFields
	}

	return []LogField{
		ErrField(err),
		ErrStackField(tracer.Trace(err)),
	}
}

func hexaErrFields(err error) []LogField {
	e := AsHexaErr(err)
	if e == nil {
		return nil
	}

	// Hexa error fields:
	fields := []LogField{
		StringField("_error_id", e.ID()),
		IntField("_http_status", e.HTTPStatus()),
	}
	for k, v := range e.Data() {
		fields = append(fields, AnyField(k, v))
	}
	for k, v := range e.ReportData() {
		fields = append(fields, AnyField(k, v))
	}

	// If exists error and error is traced,print its stack.
	fields = append(fields, ErrStackField(tracer.MoveStackIfNeeded(e, e.InternalError())))
	if e.InternalError() != nil {
		fields = append(fields, ErrField(e.InternalError()))
	}

	return fields
}

type Logger interface {

	// Core function returns the logger core concrete struct.
	// this is because sometimes we need to convert one logger
	// interface to another and need to the concrete logger.
	Core() any

	// WithCtx gets the hexa context and some keyValues
	// and return new logger contains key values as
	// log fields.
	WithCtx(ctx context.Context, args ...LogField) Logger

	// With method set key,values and return new logger
	// contains this key values as log fields.
	With(f ...LogField) Logger

	// Debug log debug message.
	Debug(msg string, args ...LogField)

	// Info log info message.
	Info(msg string, args ...LogField)

	// Message log the value as a message.
	// Use this to send message to some loggers that just want to get messages.
	// all loggers see message as info and just add simple _message tag to it.
	// but some other loggers just log messages (like our sentry logger).
	// severity of Message it just like info.
	Message(msg string, args ...LogField)

	// Warn log warning message.
	Warn(msg string, args ...LogField)

	// Error log error message
	Error(msg string, args ...LogField)
}
