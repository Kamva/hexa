package hexatranslator

import "github.com/kamva/hexa"

// emptyDriver is just empty translator driver for test purpose.
type emptyDriver struct{}

func (e emptyDriver) Localize(langs ...string) hexa.Translator {
	return e
}

func (e emptyDriver) Translate(key string, keyParams ...interface{}) (string, error) {
	return "empty_translate:" + key, nil
}

func (e emptyDriver) MustTranslate(key string, keyParams ...interface{}) string {
	t, _ := e.Translate(key, keyParams...)
	return t
}

func (e emptyDriver) TranslateDefault(key string, fallback string, keyParams ...interface{}) (string, error) {
	return e.Translate(key, keyParams...)
}

func (e emptyDriver) MustTranslateDefault(key string, fallback string, keyParams ...interface{}) string {
	t, _ := e.Translate(key, keyParams...)
	return t
}

func NewEmptyDriver() hexa.Translator {
	return emptyDriver{}
}

// Assertion
var _ hexa.Translator = &emptyDriver{}
