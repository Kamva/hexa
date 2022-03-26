package hexatranslator

import (
	"context"

	"github.com/kamva/hexa"
)

var global = NewEmptyDriver()

func SetGlobal(t hexa.Translator) {
	global = t
}

func CtxTranslator(ctx context.Context) hexa.Translator {
	if t := hexa.CtxTranslator(ctx); t != nil {
		return t
	}
	return global
}
