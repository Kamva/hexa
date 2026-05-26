package hlog

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

// recLogger records which leveled method was invoked.
type recLogger struct {
	calls *[]string
}

func (l recLogger) Core() any                                { return nil }
func (l recLogger) Enabled(Level) bool                       { return true }
func (l recLogger) WithCtx(context.Context, ...Field) Logger { return l }
func (l recLogger) With(...Field) Logger                     { return l }
func (l recLogger) Debug(string, ...Field)                   { *l.calls = append(*l.calls, "debug") }
func (l recLogger) Info(string, ...Field)                    { *l.calls = append(*l.calls, "info") }
func (l recLogger) Message(string, ...Field)                 { *l.calls = append(*l.calls, "message") }
func (l recLogger) Warn(string, ...Field)                    { *l.calls = append(*l.calls, "warn") }
func (l recLogger) Error(string, ...Field)                   { *l.calls = append(*l.calls, "error") }

func TestLevel_String(t *testing.T) {
	assert.Equal(t, "debug", DebugLevel.String())
	assert.Equal(t, "info", InfoLevel.String())
	assert.Equal(t, "warn", WarnLevel.String())
	assert.Equal(t, "error", ErrorLevel.String())
	assert.Contains(t, Level(99).String(), "99")
}

func TestLevel_CanLog(t *testing.T) {
	assert.True(t, InfoLevel.CanLog(ErrorLevel))  // info logger logs errors
	assert.True(t, InfoLevel.CanLog(InfoLevel))   // and itself
	assert.False(t, InfoLevel.CanLog(DebugLevel)) // but not debug
}

func TestLevelFromString(t *testing.T) {
	lvl, err := LevelFromString("debug")
	require.NoError(t, err)
	assert.Equal(t, DebugLevel, lvl)

	lvl, err = LevelFromString("error")
	require.NoError(t, err)
	assert.Equal(t, ErrorLevel, lvl)

	_, err = LevelFromString("nope")
	assert.Error(t, err)
}

func TestZapLevel(t *testing.T) {
	assert.Equal(t, zapcore.DebugLevel, ZapLevel(DebugLevel))
	assert.Equal(t, zapcore.ErrorLevel, ZapLevel(ErrorLevel))
	assert.Panics(t, func() { ZapLevel(Level(99)) })
}

func TestPrinterDriver(t *testing.T) {
	l := NewPrinterDriver(InfoLevel)

	assert.NotNil(t, l.Core())
	assert.False(t, l.Enabled(DebugLevel))
	assert.True(t, l.Enabled(InfoLevel))
	assert.True(t, l.Enabled(ErrorLevel))

	// With/WithCtx return a usable logger and don't mutate the original.
	assert.NotNil(t, l.With(String("k", "v")))
	assert.NotNil(t, l.WithCtx(context.Background(), Int("n", 1)))

	// Logging at/below level must not panic.
	l.Info("info msg", String("a", "b"))
	l.Debug("debug msg (suppressed)")
	l.Warn("warn")
	l.Error("err")
	l.Message("message")
}

func TestSetGlobalLogger_RebindsAndRoutesMessage(t *testing.T) {
	prev := global
	t.Cleanup(func() { SetGlobalLogger(prev) })

	var calls []string
	rl := recLogger{calls: &calls}
	SetGlobalLogger(rl)

	assert.Equal(t, rl, GlobalLogger())

	Message("m")
	Info("i")
	Warn("w")
	Error("e")
	Debug("d")

	// Message must route to the driver's Message, not Info.
	assert.Equal(t, []string{"message", "info", "warn", "error", "debug"}, calls)
}

func TestNewWriter(t *testing.T) {
	l := NewPrinterDriver(DebugLevel)

	w := NewWriter(l, InfoLevel)
	n, err := w.Write([]byte("hello"))
	require.NoError(t, err)
	assert.Equal(t, len("hello"), n)

	// nil logger falls back to the global logger.
	assert.NotNil(t, NewWriter(nil, ErrorLevel))
	// an unknown level falls back to info without panicking.
	assert.NotNil(t, NewWriter(l, Level(99)))
}

func TestMapToFields(t *testing.T) {
	fields := MapToFields(map[string]any{"a": "b"})
	require.Len(t, fields, 1)
	k, v := FieldToKeyVal(fields[0])
	assert.Equal(t, "a", k)
	assert.Equal(t, "b", v)
}

func TestErrStack(t *testing.T) {
	f := ErrStack(errors.New("boom"))
	assert.Equal(t, ErrorStackLogKey, f.Key)
}
