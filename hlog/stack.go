package hlog

import (
	"fmt"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"strings"
)

type StackedLogger interface {
	// LoggerByName returns logger by its name.
	// logger can be nil if does not exists.
	LoggerByName(name string) hexa.Logger
}

type stackedLogger struct {
	stack map[string]hexa.Logger
}

func (l *stackedLogger) LoggerByName(name string) hexa.Logger {
	return l.stack[name]
}

func (l *stackedLogger) Core() interface{} {
	return l.stack
}

func (l *stackedLogger) With(ctx hexa.Context, args ...interface{}) hexa.Logger {
	stack := make(map[string]hexa.Logger)
	for k, logger := range l.stack {
		stack[k] = logger.With(ctx, args...)
	}

	return NewStackLoggerDriverWith(stack)
}

func (l *stackedLogger) WithFields(args ...interface{}) hexa.Logger {
	stack := make(map[string]hexa.Logger)
	for k, logger := range l.stack {
		stack[k] = logger.WithFields(args...)
	}

	return NewStackLoggerDriverWith(stack)
}

func (l *stackedLogger) WithFunc(f hexa.LogFunc) hexa.Logger {
	return f(l)
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

type LoggerOptions interface {
	Zap() ZapOptions
	Sentry() SentryOptions
}

// NewStackLoggerDriver return new stacked logger .
// If logger name is invalid,it will return error.
func NewStackLoggerDriver(stackList []string, opts LoggerOptions) (hexa.Logger, error) {
	stack := make(map[string]hexa.Logger, len(stackList))

	zap := "zap"
	printer := "printer"
	sentry := "sentry"

	for _, loggerName := range stackList {
		var logger hexa.Logger
		var err error

		switch strings.ToLower(loggerName) {
		case zap:
			stack[zap] = NewZapDriver(opts.Zap())
		case printer:
			stack[printer] = NewPrinterDriver()
		case sentry:
			logger, err = NewSentryDriver(opts.Sentry())
			if err != nil {
				return nil, tracer.Trace(err)
			}
			stack[sentry] = logger
		default:
			return nil, tracer.Trace(fmt.Errorf("logger with name %s not found", loggerName))
		}
	}

	return NewStackLoggerDriverWith(stack), nil
}

// NewStackLoggerDriverWith return new instance of hexa logger with stacked logger driver.
func NewStackLoggerDriverWith(stack map[string]hexa.Logger) hexa.Logger {
	return &stackedLogger{stack}
}

// Assert stackedLogger implements hexa Logger.
var _ hexa.Logger = &stackedLogger{}
var _ StackedLogger = &stackedLogger{}
