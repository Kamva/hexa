package kittytranslator

import (
	"github.com/Kamva/gutil"
	"github.com/Kamva/kitty"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type i18nTranslator struct {
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
}

func (i i18nTranslator) Localize(langs ...string) kitty.Translator {
	return NewI18nDriver(i.bundle, i18n.NewLocalizer(i.bundle, langs...))
}

func (i i18nTranslator) translate(key string, keyParams ...interface{}) (string, error) {
	params, err := gutil.KeyValuesToMap(keyParams)

	if err != nil {
		return "", err
	}

	return i.localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: params,
	})
}

func (i i18nTranslator) Translate(key string, keyParams ...interface{}) string {
	translated, err := i.translate(key, keyParams)

	if err != nil {
		return key
	}

	return translated
}

func (i i18nTranslator) TranslateDefault(key string, fallback string, keyParams ...interface{}) string {
	translated, err := i.translate(key, keyParams)

	if err != nil {
		return fallback
	}

	return translated
}

// NewI18nDriver return new instance of i18n driver to use as kitty Translator.
func NewI18nDriver(bundle *i18n.Bundle, localizer *i18n.Localizer) kitty.Translator {
	return &i18nTranslator{bundle: bundle, localizer: localizer}
}

// Assert i18nTranslator implements kitty Translator.
var _ kitty.Translator = &i18nTranslator{}
