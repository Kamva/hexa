package hexalogger

import (
	"github.com/Kamva/hexa"
	"go.uber.org/zap"
)

type zapLogger struct {
	logger *zap.SugaredLogger
}

func (l *zapLogger) Core() interface{} {
	return l.logger
}

func (l *zapLogger) With(ctx hexa.Context, keyValues ...interface{}) hexa.Logger {
	return l.WithFields(keyValues...)
}

func (l *zapLogger) WithFields(args ...interface{}) hexa.Logger {
	if len(args) > 0 {
		return NewZapDriverWith(l.logger.With(args...))
	}
	return l
}

func (l *zapLogger) Debug(i ...interface{}) {
	l.logger.Debug(i...)
}

func (l *zapLogger) Info(i ...interface{}) {
	l.logger.Info(i...)
}

func (l *zapLogger) Message(i ...interface{}) {
	l.logger.Info(i...)
}

func (l *zapLogger) Error(i ...interface{}) {
	l.logger.Error(i...)
}

// NewZapDriver return new instance of hexa logger with zap driver.
// TODO: get log options by provided config please.
func NewZapDriver(config hexa.Config) hexa.Logger {
	l, _ := zap.NewProduction()
	return NewZapDriverWith(l.Sugar())
}

// NewZapDriver return new instance of hexa logger with zap driver.
func NewZapDriverWith(logger *zap.SugaredLogger) hexa.Logger {
	return &zapLogger{logger}
}

// Assert zapLogger implements hexa Logger.
var _ hexa.Logger = &zapLogger{}
