package hlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestStackedLogger_LoggerByName(t *testing.T) {
	o := StackOptions{
		Level:     DebugLevel,
		ZapConfig: DefaultZapConfig(true, zapcore.DebugLevel,"json"),
	}

	l, err := NewStackLoggerDriver([]string{PrinterLogger, ZapLogger}, o)

	assert.Nil(t, err)
	sl, ok := l.(StackedLogger)
	if assert.True(t, ok) {
		assert.NotNil(t, sl.LoggerByName(ZapLogger))
	}
}
