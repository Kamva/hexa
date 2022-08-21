package hexa

import (
	"context"
	"net"
	"net/http"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa/hlog"
)

// key is context key.
type contextKey string

const (
	ctxKeyRequest        contextKey = "_ctx_request"         // value MUST be *http.Request (optional)
	ctxKeyCorrelationId  contextKey = "_ctx_correlation_id"  // value MUST be cid string
	ctxKeyLocale         contextKey = "_ctx_locale"          // value MUST be locale string (can be empty string)
	ctxKeyUser           contextKey = "_ctx_user"            // Value MUST be user
	ctxKeyBaseLogger     contextKey = "_ctx_base_logger"     // value MUST be Logger interface
	ctxKeyLogger         contextKey = "_ctx_logger"          // value MUST be Logger interface
	ctxKeyBaseTranslator contextKey = "_ctx_base_translator" // value MUST be Translator interface
	ctxKeyTranslator     contextKey = "_ctx_translator"      // value MUST be Translator interface
	ctxKeyStore          contextKey = "_ctx_store"           // value MUST be Store interface.
)

func WithRequest(ctx context.Context, r *http.Request) context.Context {
	return updateLogger(context.WithValue(ctx, ctxKeyRequest, r))
}

func CtxRequest(ctx context.Context) *http.Request {
	r, _ := ctx.Value(ctxKeyRequest).(*http.Request)
	return r
}

func WithCorrelationId(ctx context.Context, cid string) context.Context {
	return updateLogger(context.WithValue(ctx, ctxKeyCorrelationId, cid))
}

func CtxCorrelationId(ctx context.Context) string {
	r, _ := ctx.Value(ctxKeyCorrelationId).(string)
	return r
}

func WithLocale(ctx context.Context, locale string) context.Context {
	return updateTranslator(context.WithValue(ctx, ctxKeyLocale, locale))
}

func CtxLocale(ctx context.Context) string {
	r, _ := ctx.Value(ctxKeyLocale).(string)
	return r
}

func WithUser(ctx context.Context, u User) context.Context {
	return updateLogger(context.WithValue(ctx, ctxKeyUser, u))
}

func CtxUser(ctx context.Context) User {
	u, _ := ctx.Value(ctxKeyUser).(User)
	return u
}

// WithBaseLogger sets the base logger in the context. when something change in the context
// we'll use this base logger to update the logger
func WithBaseLogger(ctx context.Context, l hlog.Logger) context.Context {
	return updateLogger(context.WithValue(ctx, ctxKeyBaseLogger, l))
}

func CtxBaseLogger(ctx context.Context) hlog.Logger {
	l, _ := ctx.Value(ctxKeyBaseLogger).(hlog.Logger)
	return l
}

// WithBaseTranslator sets the base Translator in the context. when you change the locale,
// we use the base translator to set the new translator.
func WithBaseTranslator(ctx context.Context, t Translator) context.Context {
	return updateTranslator(context.WithValue(ctx, ctxKeyTranslator, t))
}

func CtxBaseTranslator(ctx context.Context) Translator {
	l, _ := ctx.Value(ctxKeyBaseTranslator).(Translator)
	return l
}

func WithStore(ctx context.Context, s Store) context.Context {
	return context.WithValue(ctx, ctxKeyStore, s)
}

func CtxStore(ctx context.Context) Store {
	s, _ := ctx.Value(ctxKeyStore).(Store)
	return s
}

// WithLogger sets the final logger, this will overwrites the auto created
// logger with extra fields. If you want to set the base logger and auto
// update on every change in the context, use WithBaseLogger instead of
// WithLogger.
func WithLogger(ctx context.Context, l hlog.Logger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger, l)
}

// CtxLogger returns the context logger.
func CtxLogger(ctx context.Context) hlog.Logger {
	l, _ := ctx.Value(ctxKeyLogger).(hlog.Logger)
	return l
}

// Logger tries to get logger from the context, otherwise returns the default logger.
func Logger(ctx context.Context) hlog.Logger {
	if l := CtxLogger(ctx); l != nil {
		return l
	}
	return hlog.GlobalLogger()
}

// WithTranslator sets the final localized translator. Use WithBaseTranslator
// if you want to set the base translator and auto update the localized translator
// on every change on the context's "locale" field.
func WithTranslator(ctx context.Context, t Translator) context.Context {
	return context.WithValue(ctx, ctxKeyTranslator, t)
}

func CtxTranslator(ctx context.Context) Translator {
	t, _ := ctx.Value(ctxKeyTranslator).(Translator)
	return t
}

type ContextParams struct {
	Request        *http.Request
	CorrelationId  string
	Locale         string // Locale syntax is just same as HTTP Accept-Language header.
	User           User
	BaseLogger     hlog.Logger
	BaseTranslator Translator
	Store          Store // Optional
}

// NewContext returns new hexa Context.
func NewContext(ctx context.Context, p ContextParams) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if p.Store == nil {
		p.Store = newStore()
	}

	ctx = WithRequest(ctx, p.Request)
	ctx = WithCorrelationId(ctx, p.CorrelationId)
	ctx = WithLocale(ctx, p.Locale)
	ctx = WithUser(ctx, p.User)
	ctx = WithBaseLogger(ctx, p.BaseLogger)
	ctx = WithBaseTranslator(ctx, p.BaseTranslator)
	ctx = WithStore(ctx, p.Store)

	return ctx
}

func updateLogger(ctx context.Context) context.Context {
	if l := CtxBaseLogger(ctx); l != nil {
		return WithLogger(ctx, l.WithCtx(ctx, logFields(ctx)...))
	}

	return ctx
}

func updateTranslator(ctx context.Context) context.Context {
	t := CtxBaseTranslator(ctx)
	if t == nil {
		return ctx
	}
	if locale := CtxLocale(ctx); locale != "" {
		return WithTranslator(ctx, t.Localize(locale))
	}

	return WithTranslator(ctx, t.Localize())
}

func logFields(ctx context.Context) []hlog.Field {
	u := CtxUser(ctx)
	r := CtxRequest(ctx)
	cid := CtxCorrelationId(ctx)

	fields := make([]hlog.Field, 0)
	if u != nil {
		fields = append(fields,
			hlog.String("_user_type", string(u.Type())),
			hlog.String("_user_id", u.Identifier()),
			hlog.String("_username", u.Username()),
		)
	}
	if cid != "" {
		fields = append(fields, hlog.String("_correlation_id", cid))
	}

	if r != nil {
		rid := r.Header.Get("X-Request-ID")
		if rid != "" {
			fields = append(fields, hlog.String("_request_id", rid))
		}

		if ip, port, err := net.SplitHostPort(gutil.IP(r)); err == nil {
			fields = append(fields, hlog.String("_ip", ip))
			fields = append(fields, hlog.String("_port", port))
		}
	}
	return fields
}
