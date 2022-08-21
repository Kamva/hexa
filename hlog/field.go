package hlog

import (
	"github.com/kamva/tracer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const ErrorStackLogKey = "_stack"

// Type and function aliases from zap to limit the libraries scope into hexa code

type Field = zapcore.Field

var Int64 = zap.Int64
var Int32 = zap.Int32
var Int = zap.Int
var Uint32 = zap.Uint32
var Uint64 = zap.Uint64
var String = zap.String
var Any = zap.Any
var Err = zap.Error
var NamedErr = zap.NamedError
var Bool = zap.Bool
var Duration = zap.Duration
var Time = zap.Time
var Times = zap.Times
var Timep = zap.Timep

// ErrStack prints error stack(if exists) using logger
func ErrStack(err error) Field {
	return String(ErrorStackLogKey, tracer.StackAsString(err))
}
