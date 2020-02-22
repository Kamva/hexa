package kitty

type Logger interface {

	// Core function returns the logger core concrete struct.
	// this is because in sometimes we need to convert one logger
	// interface to another and need to logger code.
	Core()interface{}

	// With method set key,values and return new logger
	// contains this key values as log fields.
	WithFields(keyValues ...interface{}) Logger

	// Debug log debug message.
	Debug(i ...interface{})

	// Info log info message.
	Info(i ...interface{})

	// Warn log warn message.
	Warn(i ...interface{})

	// Error log error message
	Error(i ...interface{})

	// Fatal log fatal message.
	Fatal(i ...interface{})

	// Panic log message as fatal and the panic.
	Panic(i ...interface{})
}
