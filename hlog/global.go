package hlog

// initialize global with a simple printerDriver as default
// global logger until you change it in bootstrap stage of
// your app.
var global = NewPrinterDriver(DebugLevel)

func SetGlobalLogger(l Logger) {
	global = l
	Enabled = global.Enabled
	WithCtx = global.WithCtx
	With = global.With
	Debug = global.Debug
	Info = global.Info
	Message = global.Info
	Warn = global.Warn
	Error = global.Error
}

func GlobalLogger() Logger {
	return global
}

var Enabled = global.Enabled
var WithCtx = global.WithCtx
var With = global.With
var Debug = global.Debug
var Info = global.Info
var Message = global.Message
var Warn = global.Warn
var Error = global.Error
