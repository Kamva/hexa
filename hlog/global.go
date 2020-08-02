package hlog

import "github.com/kamva/hexa"

// initialize global with a simple printerDriver as default
// global logger until you change it in bootstrap stage of
// your source code.
var global = NewPrinterDriver()

func SetGlobalLogger(l hexa.Logger) {
	global = l
	With = global.With
	WithFields = global.WithFields
	Debug = global.Debug
	Info = global.Info
	Message = global.Info
	Warn = global.Warn
	Error = global.Error
}

var With = global.With
var WithFields = global.WithFields
var Debug = global.Debug
var Info = global.Info
var Message = global.Info
var Warn = global.Warn
var Error = global.Error
