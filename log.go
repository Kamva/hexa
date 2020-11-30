package hexa

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const ErrorStackLogKey = "__stack__"

// LogFunc gets logger, set data on it and return it.
// example of such functions is ErrStack function in the
// hlog package.
type LogFunc func(Logger) Logger

// Type and function aliases from zap to limit the libraries scope into hexa code
type LogField = zapcore.Field

// hlog LogField helpers to create fields from types.
var Int64Field = zap.Int64
var Int32Field = zap.Int32
var IntField = zap.Int
var Uint32Field = zap.Uint32
var StringField = zap.String
var AnyField = zap.Any
var ErrField = zap.Error
var NamedErrField = zap.NamedError
var BoolField = zap.Bool
var DurationField = zap.Duration


type Logger interface {

	// Core function returns the logger core concrete struct.
	// this is because sometimes we need to convert one logger
	// interface to another and need to the concrete logger.
	Core() interface{}

	// With get the hexa context and some keyValues
	// and return new logger contains key values as
	// log fields.
	WithCtx(ctx Context, args ...LogField) Logger

	// WithFields method set key,values and return new logger
	// contains this key values as log fields.
	With(f ...LogField) Logger

	// WithF call to the provided function to set a field in the logger using provided function.
	WithFunc(f LogFunc) Logger

	// Debug log debug message.
	Debug(msg string, args ...LogField)

	// Info log info message.
	Info(msg string, args ...LogField)

	// Message log the value as a message.
	// Use this to send message to some loggers that just want to get messages.
	// all loggers see message as info and just add simple __message__ tag to it.
	// but some other loggers just log messages (like our sentry logger).
	// severity of Message it just like info.
	Message(msg string, args ...LogField)

	// Warn log warning message.
	Warn(msg string, args ...LogField)

	// Error log error message
	Error(msg string, args ...LogField)
}
