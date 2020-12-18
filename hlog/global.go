package hlog

import "github.com/kamva/hexa"

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

var WithCtx = global.WithCtx
var With = global.With
var Debug = global.Debug
var Info = global.Info
var Message = global.Message
var Warn = global.Warn
var Error = global.Error
