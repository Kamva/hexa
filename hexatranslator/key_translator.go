package hexatranslator

import "github.com/kamva/hexa"

// keyTranslator just returns key itself as the translation result.
type keyTranslator struct{}

func (e keyTranslator) Localize(langs ...string) hexa.Translator { //nolint:revive
	return e
}

func (e keyTranslator) Translate(key string, keyParams ...any) (string, error) { //nolint:revive
	return key, nil
}

func (e keyTranslator) MustTranslate(key string, keyParams ...any) string {
	t, _ := e.Translate(key, keyParams...)
	return t
}

func (e keyTranslator) TranslateDefault(key string, fallback string, keyParams ...any) (string, error) { //nolint:revive
	return e.Translate(key, keyParams...)
}

func (e keyTranslator) MustTranslateDefault(key string, fallback string, keyParams ...any) string { //nolint:revive
	t, _ := e.Translate(key, keyParams...)
	return t
}

// NewKeyTranslator returns keyTranslator that just returns key as the translation result.
func NewKeyTranslator() hexa.Translator {
	return keyTranslator{}
}

// Assertion
var _ hexa.Translator = &keyTranslator{}
