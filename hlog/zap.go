package hlog

import (
	"context"

	"github.com/kamva/hexa"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.Logger
}

func (l *zapLogger) Core() any {
	return l.logger
}

func (l *zapLogger) WithCtx(_ context.Context, args ...Field) hexa.Logger {
	return l.With(args...)
}

func (l *zapLogger) With(args ...Field) hexa.Logger {
	if len(args) > 0 {
		return NewZapDriver(l.logger.With(args...))
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

func NewZapDriverFromConfig(cfg zap.Config) hexa.Logger {
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return NewZapDriver(l)
}

// NewZapDriver return new instance of hexa logger with zap driver.
func NewZapDriver(logger *zap.Logger) hexa.Logger {
	return &zapLogger{logger}
}

// Assert zapLogger implements hexa Logger.
var _ hexa.Logger = &zapLogger{}
