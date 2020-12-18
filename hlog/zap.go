package hlog

import (
	"github.com/kamva/hexa"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.Logger
}

func (l *zapLogger) Core() interface{} {
	return l.logger
}

func (l *zapLogger) WithCtx(ctx hexa.Context, args ...Field) hexa.Logger {
	return l.With(args...)
}

func (l *zapLogger) With(args ...Field) hexa.Logger {
	if len(args) > 0 {
		return NewZapDriverWith(l.logger.With(args...))
	}
	return l
}

func (l *zapLogger) Debug(msg string, args ...Field) {
	l.logger.Debug(msg, args...)
}

func (l *zapLogger) Info(msg string, args ...Field) {
	l.logger.Info(msg, args...)
}

func (l *zapLogger) Message(msg string, args ...Field) {
	l.logger.Info(msg, args...)
}

func (l *zapLogger) Warn(msg string, args ...Field) {
	l.logger.Warn(msg, args...)
}

func (l *zapLogger) Error(msg string, args ...Field) {
	l.logger.Error(msg, args...)
}

type ZapOptions struct {
	Debug bool
	Level zapcore.Level
}

// NewZapDriver return new instance of hexa logger with zap driver.
func NewZapDriver(o ZapOptions) hexa.Logger {
	cfg := zap.NewProductionConfig()
	if o.Debug {
		cfg = zap.NewDevelopmentConfig()
	}
	cfg.Level.SetLevel(o.Level)

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return NewZapDriverWith(l)
}

// NewZapDriver return new instance of hexa logger with zap driver.
func NewZapDriverWith(logger *zap.Logger) hexa.Logger {
	return &zapLogger{logger}
}

// Assert zapLogger implements hexa Logger.
var _ hexa.Logger = &zapLogger{}
