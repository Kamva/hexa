package hlog

import (
	"context"
)

type Logger interface {
	// Core function returns the logger core concrete struct.
	// this is because sometimes we need to convert one logger
	// interface to another and need to the concrete logger.
	Core() any

	// Enabled returns true if the logger covers the provided log level.
	// Each level covers itself and all higher logging levels.
	// For example the InfoLevel level covers WarnLevel and ErrorLevel.
	Enabled(lvl Level) bool

	// WithCtx gets the hexa context and some keyValues
	// and return new logger contains key values as
	// log fields.
	WithCtx(ctx context.Context, fields ...Field) Logger

	// With method set key,values and return new logger
	// contains this key values as log fields.
	With(f ...Field) Logger

	// Debug log debug message.
	Debug(msg string, fields ...Field)

	// Info log info message.
	Info(msg string, fields ...Field)

	// Message log the value as a message.
	// Use this to send message to some loggers that just want to get messages.
	// all loggers see message as info and just add simple _message tag to it.
	// but some other loggers just log messages (like our sentry logger).
	// severity of Message it just like info.
	Message(msg string, fields ...Field)

	// Warn log warning message.
	Warn(msg string, fields ...Field)

	// Error log error message
	Error(msg string, fields ...Field)
}
