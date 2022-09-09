package hlog

import (
	"context"
	"fmt"
	"time"
)

type printerLogger struct {
	timeFormat string
	level      Level
	with       []Field
}

func (l *printerLogger) Core() any {
	return fmt.Println
}
func (l *printerLogger) Enabled(lvl Level) bool {
	return l.level.CanLog(lvl)
}
func (l *printerLogger) cloneData() []Field {
	dst := make([]Field, len(l.with))
	copy(dst, l.with)
	return dst
}
func (l *printerLogger) clone() *printerLogger {
	return &printerLogger{
		timeFormat: l.timeFormat,
		level:      l.level,
		with:       l.cloneData(),
	}
}
func (l *printerLogger) WithCtx(_ context.Context, fields ...Field) Logger {
	clone := l.clone()
	clone.with = append(clone.with, fields...)
	return clone
}

func (l *printerLogger) With(fields ...Field) Logger {
	return l.WithCtx(context.Background(), fields...)
}

func (l *printerLogger) log(level Level, msg string, fields ...Field) {
	ll := l.With(fields...).(*printerLogger)
	t := time.Now().Format(l.timeFormat)

	if l.level.CanLog(level) {
		fmt.Println(fmt.Sprintf("%s %s: ", t, level), fieldsToMap(ll.with...), msg)
	}
}

func (l *printerLogger) Debug(msg string, fields ...Field) {
	l.log(DebugLevel, msg, fields...)
}

func (l *printerLogger) Info(msg string, fields ...Field) {
	l.log(InfoLevel, msg, fields...)
}

func (l *printerLogger) Message(msg string, fields ...Field) {
	l.log(InfoLevel, msg, fields...)
}

func (l *printerLogger) Warn(msg string, fields ...Field) {
	l.log(WarnLevel, msg, fields...)
}

func (l *printerLogger) Error(msg string, fields ...Field) {
	l.log(ErrorLevel, msg, fields...)
}

// NewPrinterDriver returns new instance of hexa logger
// with printer driver.
// Note: printer logger driver is just for test purpose.
// dont use it in production.
func NewPrinterDriver(l Level) Logger {
	return &printerLogger{
		timeFormat: "2006-01-02T15:04:05.000-0700",
		level:      l,
		with:       make([]Field, 0),
	}
}

// Assert printerLogger implements hexa Logger.
var _ Logger = &printerLogger{}
