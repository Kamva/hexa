package mgmadapter

import (
	"context"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"go.mongodb.org/mongo-driver/mongo"
)

type DBHealthChecker struct {
	client *mongo.Client
}

func (r *DBHealthChecker) HealthIdentifier() string {
	return "mongo_db"
}

func (r *DBHealthChecker) LivenessStatus(ctx context.Context) hexa.LivenessStatus {
	if err := r.client.Ping(ctx, nil); err != nil {
		hlog.Error("error on send ping to MongoDB", hlog.ErrStack(tracer.Trace(err)), hlog.Err(err))
		return hexa.StatusDead
	}
	return hexa.StatusAlive
}

func (r *DBHealthChecker) ReadinessStatus(ctx context.Context) hexa.ReadinessStatus {
	if err := r.client.Ping(ctx, nil); err != nil {
		hlog.Error("error on send ping to MongoDB", hlog.ErrStack(tracer.Trace(err)), hlog.Err(err))
		return hexa.StatusUnReady
	}
	return hexa.StatusReady
}

func NewDBHealthChecker(client *mongo.Client) *DBHealthChecker {
	return &DBHealthChecker{}
}
