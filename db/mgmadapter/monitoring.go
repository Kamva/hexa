package mgmadapter

import (
	"context"
	"strconv"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"go.mongodb.org/mongo-driver/event"
)

// composedMonitor composes multiple MongoDB server monitors.
type composedMonitor struct {
	monitors []*event.CommandMonitor
}

func NewComposedMonitor(monitors ...*event.CommandMonitor) *event.CommandMonitor {
	m := &composedMonitor{monitors: monitors}

	return &event.CommandMonitor{
		Started:   m.Started,
		Succeeded: m.Succeeded,
		Failed:    m.Failed,
	}
}

func (m *composedMonitor) Started(c context.Context, e *event.CommandStartedEvent) {
	for _, monitor := range m.monitors {
		if monitor.Started != nil {
			monitor.Started(c, e)
		}
	}
}

func (m *composedMonitor) Succeeded(c context.Context, e *event.CommandSucceededEvent) {
	for _, monitor := range m.monitors {
		if monitor.Succeeded != nil {
			monitor.Succeeded(c, e)
		}
	}
}

func (m *composedMonitor) Failed(c context.Context, e *event.CommandFailedEvent) {
	for _, monitor := range m.monitors {
		if monitor.Failed != nil {
			monitor.Failed(c, e)
		}
	}
}

// logMonitor logs every command and its reply, so be careful and
// don't use it in production mode.
type logMonitor struct {
	l hexa.Logger
}

func NewLogMonitor(l hexa.Logger) *event.CommandMonitor {
	m := &logMonitor{l: l}
	return &event.CommandMonitor{
		Started:   m.Started,
		Succeeded: m.Succeeded,
		Failed:    m.Failed,
	}
}

func (m *logMonitor) Started(c context.Context, e *event.CommandStartedEvent) {
	m.logger(c).Debug("MongoDB command started",
		hlog.String("command", e.Command.String()),
		hlog.String("db", e.DatabaseName),
		hlog.String("command_name", e.CommandName),
		hlog.String("request_id", strconv.FormatInt(e.RequestID, 10)),
		hlog.String("connection_id", e.ConnectionID),
	)
}

func (m *logMonitor) Succeeded(c context.Context, e *event.CommandSucceededEvent) {
	m.logger(c).Debug("MongoDB command succeeded",
		hlog.String("reply", e.Reply.String()),
		hlog.String("command_name", e.CommandName),
		hlog.String("request_id", strconv.FormatInt(e.RequestID, 10)),
		hlog.String("connection_id", e.ConnectionID),
	)
}

func (m *logMonitor) Failed(c context.Context, e *event.CommandFailedEvent) {
	m.logger(c).Debug("MongoDB command failed",
		hlog.String("failure", e.Failure),
		hlog.String("command_name", e.CommandName),
		hlog.String("request_id", strconv.FormatInt(e.RequestID, 10)),
		hlog.String("connection_id", e.ConnectionID),
	)
}

// logger returns the user's logger using context, otherwise
// returns the its generic logger.
func (m *logMonitor) logger(c context.Context) hexa.Logger {
	l := hexa.CtxLogger(c)
	if l == nil {
		return m.l
	}
	return l
}
