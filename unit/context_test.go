package unit

import (
	"context"
	"testing"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hexatranslator"
	"github.com/kamva/hexa/hlog"
	"go.uber.org/zap"
)

func BenchmarkNewContext(b *testing.B) {
	l := hlog.NewPrinterDriver(hlog.DebugLevel)
	t := hexatranslator.NewEmptyDriver()
	guest := hexa.NewGuest()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		hexa.NewContext(context.Background(), hexa.ContextParams{
			Request:        nil,
			CorrelationId:  "test",
			Locale:         "en-US",
			User:           guest,
			BaseLogger:     l,
			BaseTranslator: t,
		})
	}
}

func BenchmarkContextLogger(b *testing.B) {
	cfg := zap.NewProductionConfig()
	cfg.Level.SetLevel(zap.DebugLevel)

	l := hlog.NewZapDriverFromConfig(cfg)
	t := hexatranslator.NewEmptyDriver()

	ctx := hexa.NewContext(context.Background(), hexa.ContextParams{
		Request:        nil,
		CorrelationId:  "test",
		Locale:         "en-US",
		User:           hexa.NewGuest(),
		BaseLogger:     l,
		BaseTranslator: t,
	})

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		l := hlog.CtxLogger(ctx)
		_ = l
	}
}
