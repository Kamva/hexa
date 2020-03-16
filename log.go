package hexa

type Logger interface {

	// Core function returns the logger core concrete struct.
	// this is because in sometimes we need to convert one logger
	// interface to another and need to logger code.
	Core() interface{}

	// With get the hexa context and some keyValues
	// and return new logger contains key values as
	// log fields.
	With(ctx Context, keyValues ...interface{}) Logger

	// WithFields method set key,values and return new logger
	// contains this key values as log fields.
	WithFields(keyValues ...interface{}) Logger

	// Debug log debug message.
	Debug(i ...interface{})

	// Info log info message.
	Info(i ...interface{})

	// Message log the value as a message.
	// Use this to send message to some loggers that just want to get messages.
	// all loggers see message as info and just add simple __message__ tag to it.
	// but some other loggers just log messages (like our sentry logger).
	// severity of Message it just like info.
	Message(i ...interface{})

	// Error log error message
	Error(i ...interface{})
}
