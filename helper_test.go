package hexa

type emptyTranslator struct {
}

func (t emptyTranslator) Localize(langs ...string) Translator {
	return t
}

func (t emptyTranslator) Translate(key string, keyParams ...any) (string, error) {
	return "{test translate}", nil
}

func (t emptyTranslator) MustTranslate(key string, keyParams ...any) string {
	return "{test translate}"
}

func (t emptyTranslator) TranslateDefault(key string, fallback string, keyParams ...any) (string, error) {
	return "{test translate}", nil
}

func (t emptyTranslator) MustTranslateDefault(key string, fallback string, keyParams ...any) string {
	return "{test translate}"
}

var _ Translator = &emptyTranslator{}
