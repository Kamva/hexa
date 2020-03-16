package hexalogger

import (
	"github.com/Kamva/hexa"
)

type stackedLogger struct {
	stack []hexa.Logger
}

func (l *stackedLogger) Core() interface{} {
	return l.stack
}

func (l *stackedLogger) With(args ...interface{}) hexa.Logger {
	if len(args) == 0 {
		return l
	}

	stack := make([]hexa.Logger, len(l.stack))
	for i, logger := range l.stack {
		stack[i] = logger.With(args...)
	}

	return NewStackLoggerDriver(stack...)
}

func (l *stackedLogger) Debug(i ...interface{}) {
	for _, logger := range l.stack {
		logger.Debug(i...)
	}
}

func (l *stackedLogger) Info(i ...interface{}) {
	for _, logger := range l.stack {
		logger.Info(i...)
	}
}

func (l *stackedLogger) Message(i ...interface{}) {
	for _, logger := range l.stack {
		logger.Message(i...)
	}
}

func (l *stackedLogger) Error(i ...interface{}) {
	for _, logger := range l.stack {
		logger.Error(i...)
	}
}

// NewStackLoggerDriver return new instance of hexa logger with stacked logger driver.
func NewStackLoggerDriver(stack ...hexa.Logger) hexa.Logger {
	return &stackedLogger{stack}
}

// Assert stackedLogger implements hexa Logger.
var _ hexa.Logger = &stackedLogger{}
