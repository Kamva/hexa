package kitty

type (
	Translator interface {
		// Localize function returns new localized translator function.
		Localize(langs ...string) Translator

		// Translate get key and params nad return translation.
		Translate(key string, keyParams ...interface{}) (string, error)

		// MustTranslate get key and params and translate,
		//otherwise panic relative error.
		MustTranslate(key string, keyParams ...interface{}) string

		// TranslateDefault translate with default message.
		TranslateDefault(key string, fallback string, keyParams ...interface{}) (string, error)

		// MustTranslateDefault translate with default message, on occur error,will panic it.
		MustTranslateDefault(key string, fallback string, keyParams ...interface{}) string
	}
)

// TranslateKeyEmptyMessage is special key that translators return empty string for that.
var TranslateKeyEmptyMessage = "__empty_message__"
