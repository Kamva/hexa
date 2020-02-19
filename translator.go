package kitty

type Translator interface {
	// Localize function returns new localized translator function.
	Localize(langs ...string) Translator

	// Translate method get key and keyParams to translate.
	// it can return key if translation not found.
	Translate(key string, keyParams ...interface{}) string

	TranslateDefault(key string, fallback string, keyParams ...interface{}) string
}
