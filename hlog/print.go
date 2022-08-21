package hlog

import (
	"context"
	"fmt"
	"time"
)

// TODO: fix bugs for duration(shows duration as nil) and nil error(does not show the "err" as key in map props) :
// e.g. (use printer driver instead of zp),
// z, err := zap.NewDevelopment()
//	gutil.PanicErr(err)
//	l := hlog.NewZapDriver(z)
//	l.Info("salam",
//		hlog.Err(errors.New("errrr")),
//		hlog.String("sa", "lam"),
//		hlog.Duration("i", time.Second),
//	)
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
	for i, v := range l.with {
		dst[i] = v
	}
	return dst
}
func (l *printerLogger) clone() *printerLogger {
	return &printerLogger{
		timeFormat: l.timeFormat,
		level:      l.level,
		with:       l.cloneData(),
	}
}
func (l *printerLogger) WithCtx(_ context.Context, args ...Field) Logger {
	clone := l.clone()
	clone.with = append(clone.with, args...)
	return clone
}

func (l *printerLogger) With(args ...Field) Logger {
	return l.WithCtx(nil, args...)
}

func (l *printerLogger) log(level Level, msg string, args ...Field) {
	ll := l.With(args...).(*printerLogger)
	t := time.Now().Format(l.timeFormat)

	if l.level.CanLog(level) {
		fmt.Println(fmt.Sprintf("%s %s: ", t, level), fieldsToMap(ll.with...), msg)
	}
}

func (l *printerLogger) Debug(msg string, args ...Field) {
	l.log(DebugLevel, msg, args...)
}

func (l *printerLogger) Info(msg string, args ...Field) {
	l.log(InfoLevel, msg, args...)
}

func (l *printerLogger) Message(msg string, args ...Field) {
	l.log(InfoLevel, msg, args...)
}

func (l *printerLogger) Warn(msg string, args ...Field) {
	l.log(WarnLevel, msg, args...)
}

func (l *printerLogger) Error(msg string, args ...Field) {
	l.log(ErrorLevel, msg, args...)
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
