package kittytranslator

import (
	"github.com/Kamva/gutil"
	"github.com/Kamva/kitty"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type i18nTranslator struct {
	fallbackLangs []string
	bundle        *i18n.Bundle
	localizer     *i18n.Localizer
}

func (i i18nTranslator) Localize(langs ...string) kitty.Translator {
	langs = append(langs, i.fallbackLangs...)

	return NewI18nDriver(i.bundle, i18n.NewLocalizer(i.bundle, langs...), i.fallbackLangs)
}

func (i i18nTranslator) Translate(key string, keyParams ...interface{}) (string, error) {
	params, err := gutil.KeyValuesToMap(keyParams)

	if err != nil {
		return "", err
	}

	return i.localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: params,
	})
}

func (i i18nTranslator) MustTranslate(key string, keyParams ...interface{}) string {
	msg, err := i.Translate(key, keyParams...)

	if err != nil {
		panic(err)
	}

	return msg
}

func (i i18nTranslator) TranslateDefault(key string, fallback string, keyParams ...interface{}) (string, error) {
	params, err := gutil.KeyValuesToMap(keyParams)

	if err != nil {
		return "", err
	}

	return i.localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:  key,
			One: fallback,
		},
		TemplateData: params,
	})
}

func (i i18nTranslator) MustTranslateDefault(key string, fallback string, keyParams ...interface{}) string {
	msg, err := i.TranslateDefault(key, fallback, keyParams...)

	if err != nil {
		panic(err)
	}

	return msg
}

// NewI18nDriver return new instance of i18n driver to use as kitty Translator.
func NewI18nDriver(bundle *i18n.Bundle, localizer *i18n.Localizer, fallbackLangs []string) kitty.Translator {
	return &i18nTranslator{bundle: bundle, localizer: localizer, fallbackLangs: fallbackLangs}
}

// Assert i18nTranslator implements kitty Translator.
var _ kitty.Translator = &i18nTranslator{}
