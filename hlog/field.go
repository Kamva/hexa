package hlog

import (
	"github.com/kamva/hexa"
)

type Field = hexa.LogField

// hlog LogField helpers to create fields from types.
var Int64 = hexa.Int64Field
var Int32 = hexa.Int32Field
var Int = hexa.IntField
var Uint32 = hexa.Uint32Field
var Uint64 = hexa.Uint64Field
var String = hexa.StringField
var Any = hexa.AnyField
var Err = hexa.ErrField
var NamedErr = hexa.NamedErrField
var Bool = hexa.BoolField
var Duration = hexa.DurationField
var Time = hexa.TimeField
var Times = hexa.TimesField
var Timep = hexa.TimepField
var ErrStack = hexa.ErrStackField
var ErrFields = hexa.ErrFields
