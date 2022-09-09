package logdriver

import (
	"context"

	"github.com/kamva/hexa/hlog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.Logger
}

func (l *zapLogger) Core() any {
	return l.logger
}

func (l *zapLogger) Enabled(lvl hlog.Level) bool {
	return l.logger.Core().Enabled(hlog.ZapLevel(lvl))
}

func (l *zapLogger) WithCtx(_ context.Context, fields ...hlog.Field) hlog.Logger {
	return l.With(fields...)
}

func (l *zapLogger) With(fields ...hlog.Field) hlog.Logger {
	if len(fields) > 0 {
		return NewZapDriver(l.logger.With(fields...))
	}
	return l
}

func (l *zapLogger) Debug(msg string, fields ...hlog.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *zapLogger) Info(msg string, fields ...hlog.Field) {
	l.logger.Info(msg, fields...)
}

func (l *zapLogger) Message(msg string, fields ...hlog.Field) {
	l.logger.Info(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...hlog.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...hlog.Field) {
	l.logger.Error(msg, fields...)
}

type ZapOptions struct {
	Debug bool
	Level zapcore.Level
}

// DefaultZapConfig generate zap config using default values.
// You can leave encoding empty to set to the default value
// which is json.
func DefaultZapConfig(debug bool, level zapcore.Level, encoding string) zap.Config {
	if encoding == "" {
		encoding = "json"
	}

	cfg := zap.NewProductionConfig()
	if debug {
		cfg = zap.NewDevelopmentConfig()
	}

	cfg.Level.SetLevel(level)
	cfg.Encoding = encoding

	return cfg
}

func NewZapDriverFromConfig(cfg zap.Config) hlog.Logger {
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return NewZapDriver(l)
}

// NewZapDriver return new instance of hexa logger with zap driver.
func NewZapDriver(logger *zap.Logger) hlog.Logger {
	return &zapLogger{logger}
}

// Assert zapLogger implements hexa Logger.
var _ hlog.Logger = &zapLogger{}
