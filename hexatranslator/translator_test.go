package hexatranslator

import (
	"context"
	"testing"

	"github.com/kamva/hexa"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestEmptyDriver(t *testing.T) {
	d := NewEmptyDriver()

	assert.Equal(t, d, d.Localize("en"))

	got, err := d.Translate("foo")
	require.NoError(t, err)
	assert.Equal(t, "empty_translate:foo", got)
	assert.Equal(t, "empty_translate:foo", d.MustTranslate("foo"))

	got, err = d.TranslateDefault("foo", "fallback")
	require.NoError(t, err)
	assert.Equal(t, "empty_translate:foo", got)
	assert.Equal(t, "empty_translate:foo", d.MustTranslateDefault("foo", "fallback"))
}

func TestKeyTranslator(t *testing.T) {
	d := NewKeyTranslator()

	assert.Equal(t, d, d.Localize())

	got, err := d.Translate("some.key")
	require.NoError(t, err)
	assert.Equal(t, "some.key", got)
	assert.Equal(t, "some.key", d.MustTranslate("some.key"))
	assert.Equal(t, "some.key", d.MustTranslateDefault("some.key", "fallback"))
}

func TestGlobal(t *testing.T) {
	prev := global
	t.Cleanup(func() { global = prev })

	SetGlobal(NewKeyTranslator())

	// With no translator in the context, the global is returned.
	assert.Equal(t, "x", CtxTranslator(context.Background()).MustTranslate("x"))

	// A translator in the context takes precedence over the global.
	ctx := hexa.WithTranslator(context.Background(), NewEmptyDriver())
	assert.Equal(t, "empty_translate:x", CtxTranslator(ctx).MustTranslate("x"))
}

func i18nDriver(t *testing.T) hexa.Translator {
	t.Helper()
	bundle := i18n.NewBundle(language.English)
	require.NoError(t, bundle.AddMessages(language.English,
		&i18n.Message{ID: "hello", Other: "Hello {{.Name}}"},
	))
	return NewI18nDriver(bundle, i18n.NewLocalizer(bundle, "en"), []string{"en"})
}

func TestI18n_Translate(t *testing.T) {
	tr := i18nDriver(t)
	got, err := tr.Translate("hello", "Name", "World")
	require.NoError(t, err)
	assert.Equal(t, "Hello World", got)
}

func TestI18n_Localize(t *testing.T) {
	got, err := i18nDriver(t).Localize("en").Translate("hello", "Name", "Go")
	require.NoError(t, err)
	assert.Equal(t, "Hello Go", got)
}

func TestI18n_EmptyMessageKey(t *testing.T) {
	got, err := i18nDriver(t).Translate(hexa.TranslateKeyEmptyMessage)
	require.NoError(t, err)
	assert.Equal(t, "", got)
}

func TestI18n_TranslateDefaultUsesFallbackForMissingKey(t *testing.T) {
	got, err := i18nDriver(t).TranslateDefault("missing.key", "the fallback")
	require.NoError(t, err)
	assert.Equal(t, "the fallback", got)
}

func TestI18n_TranslateMissingKeyReturnsError(t *testing.T) {
	_, err := i18nDriver(t).Translate("missing.key")
	assert.Error(t, err)
}
