package hexalogger

import (
	"fmt"
	"github.com/Kamva/hexa"
	"github.com/Kamva/tracer"
)

type stackedLogger struct {
	stack []hexa.Logger
}

var (
	LogConfigKeyStack = "log.stack"
)

func (l *stackedLogger) Core() interface{} {
	return l.stack
}

func (l *stackedLogger) With(ctx hexa.Context, args ...interface{}) hexa.Logger {
	stack := make([]hexa.Logger, len(l.stack))
	for i, logger := range l.stack {
		stack[i] = logger.With(ctx, args...)
	}

	return NewStackLoggerDriverWith(stack...)
}

func (l *stackedLogger) WithFields(args ...interface{}) hexa.Logger {
	stack := make([]hexa.Logger, len(l.stack))
	for i, logger := range l.stack {
		stack[i] = logger.WithFields(args...)
	}

	return NewStackLoggerDriverWith(stack...)
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

func (l *stackedLogger) Warn(i ...interface{}) {
	for _, logger := range l.stack {
		logger.Warn(i...)
	}
}

func (l *stackedLogger) Error(i ...interface{}) {
	for _, logger := range l.stack {
		logger.Error(i...)
	}
}

// NewStackLoggerDriver return new stacked logger .
// If logger name is invalid,it will return error.
func NewStackLoggerDriver(cfg hexa.Config) (hexa.Logger, error) {
	stackList := cfg.GetList(LogConfigKeyStack)
	stack := make([]hexa.Logger, len(stackList))

	for i, loggerName := range stackList {
		var logger hexa.Logger
		var err error

		switch loggerName {
		case "zap":
			logger = NewZapDriver(cfg)
		case "sentry":
			logger, err = NewSentryDriver(cfg)
			if err != nil {
				return nil, tracer.Trace(err)
			}
		default:
			return nil, tracer.Trace(fmt.Errorf("logger with name %s not found", loggerName))
		}

		stack[i] = logger
	}

	return NewStackLoggerDriverWith(stack...), nil
}

// NewStackLoggerDriverWith return new instance of hexa logger with stacked logger driver.
func NewStackLoggerDriverWith(stack ...hexa.Logger) hexa.Logger {
	return &stackedLogger{stack}
}

// Assert stackedLogger implements hexa Logger.
var _ hexa.Logger = &stackedLogger{}
