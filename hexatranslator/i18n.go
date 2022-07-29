package hexatranslator

import (
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type i18nTranslator struct {
	fallbackLangs []string
	bundle        *i18n.Bundle
	localizer     *i18n.Localizer
}

func (t i18nTranslator) Localize(langs ...string) hexa.Translator {
	langs = append(langs, t.fallbackLangs...)

	return NewI18nDriver(t.bundle, i18n.NewLocalizer(t.bundle, langs...), t.fallbackLangs)
}

func (t i18nTranslator) isEmptyMessageKey(key string) bool {
	return key == hexa.TranslateKeyEmptyMessage
}

func (t i18nTranslator) Translate(key string, keyParams ...any) (string, error) {

	if t.isEmptyMessageKey(key) {
		return "", nil
	}

	params, err := gutil.KeyValuesToMap(keyParams...)

	if err != nil {
		return "", tracer.Trace(err)
	}

	return t.localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: params,
	})
}

func (t i18nTranslator) MustTranslate(key string, keyParams ...any) string {

	msg, err := t.Translate(key, keyParams...)

	if err != nil {
		panic(err)
	}

	return msg
}

func (t i18nTranslator) TranslateDefault(key string, fallback string, keyParams ...any) (string, error) {
	if t.isEmptyMessageKey(key) {
		return "", nil
	}

	params, err := gutil.KeyValuesToMap(keyParams...)

	if err != nil {
		return "", tracer.Trace(err)
	}

	return t.localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    key,
			Zero:  fallback,
			One:   fallback,
			Two:   fallback,
			Few:   fallback,
			Many:  fallback,
			Other: fallback,
		},
		TemplateData: params,
	})
}

func (t i18nTranslator) MustTranslateDefault(key string, fallback string, keyParams ...any) string {
	msg, err := t.TranslateDefault(key, fallback, keyParams...)

	if err != nil {
		panic(err)
	}

	return msg
}

// NewI18nDriver return new instance of i18n driver to use as hexa Translator.
func NewI18nDriver(bundle *i18n.Bundle, localizer *i18n.Localizer, fallbackLangs []string) hexa.Translator {
	return &i18nTranslator{bundle: bundle, localizer: localizer, fallbackLangs: fallbackLangs}
}

// Assert i18nTranslator implements hexa Translator.
var _ hexa.Translator = &i18nTranslator{}
