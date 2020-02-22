package kitty

type Context interface {
	User() User
	Translator() Translator
	Log() Logger
}
