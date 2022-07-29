package hexa

type (
	Translator interface {
		// Localize function returns new localized translator function.
		Localize(langs ...string) Translator

		// Translate get key and params nad return translation.
		Translate(key string, keyParams ...any) (string, error)

		// MustTranslate get key and params and translate,
		//otherwise panic relative error.
		MustTranslate(key string, keyParams ...any) string

		// TranslateDefault translate with default message.
		TranslateDefault(key string, fallback string, keyParams ...any) (string, error)

		// MustTranslateDefault translate with default message, on occur error,will panic it.
		MustTranslateDefault(key string, fallback string, keyParams ...any) string
	}
)

// TranslateKeyEmptyMessage is special key that translators return empty string for that.
const TranslateKeyEmptyMessage = "_empty_message"
