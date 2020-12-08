package hlog

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
)

type SentryOptions struct {
	DSN         string
	Debug       bool
	Environment string
}

type sentryLogger struct {
	hub *sentry.Hub
}

func (l *sentryLogger) Core() interface{} {
	return l.hub
}

func (l *sentryLogger) addArgsToScope(scope *sentry.Scope, args []Field) {
	if len(args) == 0 {
		return
	}
	fields := fieldsToMap(args...)
	for key, val := range fields {
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

	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		u.IPAddress = ip
	}

	scope.SetUser(u)
}

func (l *sentryLogger) WithCtx(ctx hexa.Context, args ...Field) hexa.Logger {
	hub := l.hub.Clone()
	scope := hub.Scope()

	r := ctx.Request()
	if r != nil {
		scope.SetRequest(r)
	}

	l.setUser(scope, ctx.User(), r)

	l.addArgsToScope(scope, args)
	return NewSentryDriverWith(hub)
}

// With get some fields and set check if field's key start and end
// with single '_' character, then insert it as tag, otherwise
// insert it as extra data.
func (l *sentryLogger) With(args ...Field) hexa.Logger {
	hub := l.hub.Clone()
	l.addArgsToScope(hub.Scope(), args)
	return NewSentryDriverWith(hub)
}

func (l *sentryLogger) WithFunc(f hexa.LogFunc) hexa.Logger {
	return f(l)
}

func (l *sentryLogger) Debug(msg string, args ...Field) {
	// For now we do not capture debug messages in sentry.
}

func (l *sentryLogger) Info(msg string, args ...Field) {
	// For now we do not capture messages in info .
}

func (l *sentryLogger) Message(msg string, args ...Field) {
	l.With(args...).(*sentryLogger).hub.CaptureMessage(msg)
}

func (l *sentryLogger) Warn(msg string, args ...Field) {
	// For now we do not capture message in warn.
}

func (l *sentryLogger) Error(msg string, args ...Field) {
	l.With(args...).(*sentryLogger).hub.CaptureException(errors.New(msg))
}

// NewSentryDriver return new instance of hexa logger with sentry driver.
func NewSentryDriver(o SentryOptions) (hexa.Logger, error) {
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
func NewSentryDriverWith(hub *sentry.Hub) hexa.Logger {
	return &sentryLogger{hub}
}

// Assert sentryLogger implements hexa Logger.
var _ hexa.Logger = &sentryLogger{}
