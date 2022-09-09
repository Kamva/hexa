package hexa

import (
	"context"

	"github.com/kamva/hexa/hlog"
)

type Ping func(ctx context.Context) error

type pingHealth struct {
	l          hlog.Logger
	identifier string
	ping       Ping
	tags       map[string]string
}

func NewPingHealth(l hlog.Logger, identifier string, ping Ping, tags map[string]string) Health {
	return &pingHealth{
		l:          l,
		identifier: identifier,
		ping:       ping,
		tags:       tags,
	}
}

func (h *pingHealth) HealthIdentifier() string {
	return h.identifier
}

func (h *pingHealth) LivenessStatus(ctx context.Context) LivenessStatus {
	if err := h.ping(ctx); err != nil {
		h.l.Error("can not ping", hlog.String("health_identifier", h.identifier), hlog.Err(err))
		return StatusDead
	}

	return StatusAlive
}

func (h *pingHealth) ReadinessStatus(ctx context.Context) ReadinessStatus {
	if h.LivenessStatus(ctx) == StatusAlive {
		return StatusReady
	}

	return StatusUnReady
}

func (h *pingHealth) HealthStatus(ctx context.Context) HealthStatus {
	liveness := h.LivenessStatus(ctx)
	readiness := StatusReady
	if liveness == StatusDead {
		readiness = StatusUnReady
	}

	return HealthStatus{
		Id:    h.HealthIdentifier(),
		Alive: liveness,
		Ready: readiness,
		Tags:  h.tags,
	}
}

var _ Health = &pingHealth{}
