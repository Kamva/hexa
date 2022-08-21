package logdriver

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
)

type SentryOptions struct {
	DSN         string
	Debug       bool
	Environment string
}

type sentryLogger struct {
	hub *sentry.Hub
}

func (l *sentryLogger) Core() any {
	return l.hub
}

func (l *sentryLogger) Enabled(lvl hlog.Level) bool {
	// For the Sentry we just support error level.
	// Currently, we don't have a way to tell log is enabled
	// for the Message() method on sentry or not. but we
	// don't need to it right now.
	return lvl == hlog.ErrorLevel
}

func (l *sentryLogger) addFieldsToScope(scope *sentry.Scope, fields []hlog.Field) {
	if len(fields) == 0 {
		return
	}
	for _, arg := range fields {
		key, val := hlog.FieldToKeyVal(arg)

		// Just keys that begin and end with "_", set as tags.
		if len(key) >= 2 && key[0] == '_' && key[len(key)-1] == '_' {
			scope.SetTag(key, fmt.Sprintf("%v", val))
		} else {
			scope.SetExtra(key, val)
		}
	}
}

func (l *sentryLogger) setUser(scope *sentry.Scope, user hexa.User, r *http.Request) {
	u := sentry.User{
		IPAddress: gutil.IP(r),
		Email:     user.Email(),
		ID:        user.Identifier(),
		Username:  user.Username(),
	}

	//if ip, _, err := net.SplitHostPort(gutil.IP(r)); err == nil {
	//	u.IPAddress = ip
	//}

	scope.SetUser(u)
}

func (l *sentryLogger) WithCtx(ctx context.Context, args ...hlog.Field) hlog.Logger {
	hub := l.hub.Clone()
	scope := hub.Scope()

	r := hexa.CtxRequest(ctx)
	if r != nil {
		scope.SetRequest(r)
	}

	if user := hexa.CtxUser(ctx); user != nil {
		l.setUser(scope, user, r)
	}

	l.addFieldsToScope(scope, args)
	return NewSentryDriverWith(hub)
}

// With get some fields and set check if field's key start and end
// with single '_' character, then insert it as tag, otherwise
// insert it as extra data.
func (l *sentryLogger) With(args ...hlog.Field) hlog.Logger {
	hub := l.hub.Clone()
	l.addFieldsToScope(hub.Scope(), args)
	return NewSentryDriverWith(hub)
}

func (l *sentryLogger) Debug(msg string, args ...hlog.Field) {
	// For now we do not capture debug messages in sentry.
}

func (l *sentryLogger) Info(msg string, args ...hlog.Field) {
	// For now we do not capture messages in info .
}

func (l *sentryLogger) Message(msg string, args ...hlog.Field) {
	l.With(args...).(*sentryLogger).hub.CaptureMessage(msg)
}

func (l *sentryLogger) Warn(msg string, args ...hlog.Field) {
	// For now we do not capture message in warn.
}

func (l *sentryLogger) Error(msg string, args ...hlog.Field) {
	l.With(args...).(*sentryLogger).hub.CaptureException(errors.New(msg))
}

// NewSentryDriver return new instance of hexa logger with sentry driver.
func NewSentryDriver(o SentryOptions) (hlog.Logger, error) {
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:         o.DSN,
		Debug:       o.Debug,
		Environment: o.Environment,
	})
	if err != nil {
		return nil, err
	}
	return NewSentryDriverWith(sentry.NewHub(client, sentry.NewScope())), nil
}

// NewSentryDriverWith get the sentry hub and returns new instance
//of sentry driver for hexa logger.
func NewSentryDriverWith(hub *sentry.Hub) hlog.Logger {
	return &sentryLogger{hub}
}

// Assert sentryLogger implements hexa Logger.
var _ hlog.Logger = &sentryLogger{}
