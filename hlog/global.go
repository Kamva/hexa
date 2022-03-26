package hlog

import (
	"context"

	"github.com/kamva/hexa"
)

// initialize global with a simple printerDriver as default
// global logger until you change it in bootstrap stage of
// your app.
var global = NewPrinterDriver(DebugLevel)

func SetGlobalLogger(l hexa.Logger) {
	global = l
	WithCtx = global.WithCtx
	With = global.With
	Debug = global.Debug
	Info = global.Info
	Message = global.Info
	Warn = global.Warn
	Error = global.Error
}

// CtxLogger returns the context logger with fall back to
// the global logger
func CtxLogger(ctx context.Context) hexa.Logger {
	if l := hexa.CtxLogger(ctx); l != nil {
		return l
	}
	return global
}

var WithCtx = global.WithCtx
var With = global.With
var Debug = global.Debug
var Info = global.Info
var Message = global.Message
var Warn = global.Warn
var Error = global.Error
