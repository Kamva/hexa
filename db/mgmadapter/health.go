package mgmadapter

import (
	"context"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"go.mongodb.org/mongo-driver/mongo"
)

type DBHealth struct {
	client *mongo.Client
}

func (h *DBHealth) HealthIdentifier() string {
	return "mongo_db"
}

func (h *DBHealth) LivenessStatus(ctx context.Context) hexa.LivenessStatus {
	if err := h.client.Ping(ctx, nil); err != nil {
		hlog.Error("error on send ping to MongoDB", hlog.ErrStack(tracer.Trace(err)), hlog.Err(err))
		return hexa.StatusDead
	}
	return hexa.StatusAlive
}

func (h *DBHealth) ReadinessStatus(ctx context.Context) hexa.ReadinessStatus {
	if err := h.client.Ping(ctx, nil); err != nil {
		hlog.Error("error on send ping to MongoDB", hlog.ErrStack(tracer.Trace(err)), hlog.Err(err))
		return hexa.StatusUnReady
	}
	return hexa.StatusReady
}

func (h *DBHealth) HealthStatus(ctx context.Context) hexa.HealthStatus {
	return hexa.HealthStatus{
		Id:    h.HealthIdentifier(),
		Alive: h.LivenessStatus(ctx),
		Ready: h.ReadinessStatus(ctx),
	}
}

func NewDBHealth(client *mongo.Client) hexa.Health {
	return &DBHealth{}
}

var _ hexa.Health = &DBHealth{}
