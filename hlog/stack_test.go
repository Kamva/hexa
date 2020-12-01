package hlog

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestStackedLogger_LoggerByName(t *testing.T) {
	o := StackOptions{
		Level: DebugLevel,
		ZapOpts: &ZapOptions{
			Debug: true,
			Level: zapcore.DebugLevel,
		},
	}

	l, err := NewStackLoggerDriver([]string{PrinterLogger, ZapLogger}, o)

	assert.Nil(t, err)
	sl, ok := l.(StackedLogger)
	if assert.True(t, ok) {
		assert.NotNil(t, sl.LoggerByName(ZapLogger))
	}
}
