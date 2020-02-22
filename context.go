package kitty

import "context"

type Context interface {
	context.Context

	User() User
	Translator() Translator
	Log() Logger
}