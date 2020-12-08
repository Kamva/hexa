package hexa

type emptyLogger struct {
}

type emptyTranslator struct {
}

func (e emptyLogger) Core() interface{} {
	return nil
}

func (e emptyLogger) WithCtx(ctx Context, args ...LogField) Logger {
	return e
}

func (e emptyLogger) With(f ...LogField) Logger {
	return e
}

func (e emptyLogger) WithFunc(f LogFunc) Logger {
	return e
}

func (e emptyLogger) Debug(msg string, args ...LogField) {
}

func (e emptyLogger) Info(msg string, args ...LogField) {
}

func (e emptyLogger) Message(msg string, args ...LogField) {
}

func (e emptyLogger) Warn(msg string, args ...LogField) {
}

func (e emptyLogger) Error(msg string, args ...LogField) {
}

func (t emptyTranslator) Localize(langs ...string) Translator {
	return t
}

func (t emptyTranslator) Translate(key string, keyParams ...interface{}) (string, error) {
	return "{test translate}", nil
}

func (t emptyTranslator) MustTranslate(key string, keyParams ...interface{}) string {
	return "{test translate}"
}

func (t emptyTranslator) TranslateDefault(key string, fallback string, keyParams ...interface{}) (string, error) {
	return "{test translate}", nil
}

func (t emptyTranslator) MustTranslateDefault(key string, fallback string, keyParams ...interface{}) string {
	return "{test translate}"
}

var _ Logger = &emptyLogger{}
var _ Translator = &emptyTranslator{}
