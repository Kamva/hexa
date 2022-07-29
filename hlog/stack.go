package hlog

import (
	"context"
	"fmt"
	"strings"

	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"go.uber.org/zap"
)

const (
	ZapLogger     = "zap"
	SentryLogger  = "sentry"
	PrinterLogger = "printer"
)

type StackedLogger interface {
	// LoggerByName returns logger by its name.
	// logger can be nil if does not exists.
	LoggerByName(name string) hexa.Logger

	hexa.Bootable
	hexa.Runnable
	hexa.Shutdownable
}

type stackedLogger struct {
	stack map[string]hexa.Logger
}

func (l *stackedLogger) LoggerByName(name string) hexa.Logger {
	return l.stack[name]
}

func (l *stackedLogger) Core() any {
	return l.stack
}

func (l *stackedLogger) WithCtx(ctx context.Context, args ...Field) hexa.Logger {
	stack := make(map[string]hexa.Logger)
	for k, logger := range l.stack {
		stack[k] = logger.WithCtx(ctx, args...)
	}

	return NewStackLoggerDriverWith(stack)
}

func (l *stackedLogger) With(args ...Field) hexa.Logger {
	stack := make(map[string]hexa.Logger)
	for k, logger := range l.stack {
		stack[k] = logger.With(args...)
	}

	return NewStackLoggerDriverWith(stack)
}

func (l *stackedLogger) Debug(msg string, args ...Field) {
	for _, logger := range l.stack {
		logger.Debug(msg, args...)
	}
}

func (l *stackedLogger) Info(msg string, args ...Field) {
	for _, logger := range l.stack {
		logger.Info(msg, args...)
	}
}

func (l *stackedLogger) Message(msg string, args ...Field) {
	for _, logger := range l.stack {
		logger.Message(msg, args...)
	}
}

func (l *stackedLogger) Warn(msg string, args ...Field) {
	for _, logger := range l.stack {
		logger.Warn(msg, args...)
	}
}

func (l *stackedLogger) Error(msg string, args ...Field) {
	for _, logger := range l.stack {
		logger.Error(msg, args...)
	}
}

func (l *stackedLogger) Boot() error {
	for _, logger := range l.stack {
		if bootable, ok := logger.(hexa.Bootable); ok {
			if err := bootable.Boot(); err != nil {
				return tracer.Trace(err)
			}
		}
	}

	return nil
}

func (l *stackedLogger) Run() error {
	for _, logger := range l.stack {
		if runnable, ok := logger.(hexa.Runnable); ok {
			if err := runnable.Run(); err != nil {
				return tracer.Trace(err)
			}
		}
	}

	return nil
}

func (l *stackedLogger) Shutdown(ctx context.Context) error {
	for _, logger := range l.stack {
		if runnable, ok := logger.(hexa.Shutdownable); ok {
			if err := runnable.Shutdown(ctx); err != nil {
				return tracer.Trace(err)
			}

			// If ctx is closed
			if err := ctx.Err(); err != nil {
				return tracer.Trace(err)
			}
		}
	}

	return nil
}

type StackOptions struct {
	Level      Level
	ZapConfig  zap.Config
	SentryOpts *SentryOptions
}

// NewStackLoggerDriver return new stacked logger .
// If logger name is invalid,it will return error.
func NewStackLoggerDriver(stackList []string, opts StackOptions) (hexa.Logger, error) {
	stack := make(map[string]hexa.Logger, len(stackList))

	for _, loggerName := range stackList {
		var logger hexa.Logger
		var err error

		switch strings.ToLower(loggerName) {
		case ZapLogger:
			stack[ZapLogger] = NewZapDriverFromConfig(opts.ZapConfig)
		case PrinterLogger:
			stack[PrinterLogger] = NewPrinterDriver(opts.Level)
		case SentryLogger:
			logger, err = NewSentryDriver(*opts.SentryOpts)
			if err != nil {
				return nil, tracer.Trace(err)
			}
			stack[SentryLogger] = logger
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
