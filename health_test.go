package hexa

import (
	"context"
	"errors"
	"testing"

	"github.com/kamva/hexa/hlog"
	"github.com/stretchr/testify/assert"
)

func aliveHealth(id string) Health {
	return NewPingHealth(hlog.NewPrinterDriver(hlog.ErrorLevel), id, func(context.Context) error { return nil }, nil)
}

func deadHealth(id string) Health {
	return NewPingHealth(hlog.NewPrinterDriver(hlog.ErrorLevel), id, func(context.Context) error { return errors.New("down") }, nil)
}

func TestPingHealth(t *testing.T) {
	ctx := context.Background()

	ok := NewPingHealth(hlog.NewPrinterDriver(hlog.ErrorLevel), "id1",
		func(context.Context) error { return nil }, map[string]string{"k": "v"})
	assert.Equal(t, "id1", ok.HealthIdentifier())
	assert.Equal(t, StatusAlive, ok.LivenessStatus(ctx))
	assert.Equal(t, StatusReady, ok.ReadinessStatus(ctx))
	hs := ok.HealthStatus(ctx)
	assert.Equal(t, StatusAlive, hs.Alive)
	assert.Equal(t, StatusReady, hs.Ready)
	assert.Equal(t, "v", hs.Tags["k"])

	bad := deadHealth("id2")
	assert.Equal(t, StatusDead, bad.LivenessStatus(ctx))
	assert.Equal(t, StatusUnReady, bad.ReadinessStatus(ctx))
	assert.Equal(t, StatusUnReady, bad.HealthStatus(ctx).Ready)
}

func TestHealthReporter_AllHealthy(t *testing.T) {
	ctx := context.Background()
	r := NewHealthReporter().AddToChecks(aliveHealth("a"))

	assert.Equal(t, StatusAlive, r.LivenessStatus(ctx))
	assert.Equal(t, StatusReady, r.ReadinessStatus(ctx))

	report := r.HealthReport(ctx)
	assert.Equal(t, StatusAlive, report.Alive)
	assert.Equal(t, StatusReady, report.Ready)
	assert.Len(t, report.Statuses, 1)
}

func TestHealthReporter_Unhealthy(t *testing.T) {
	ctx := context.Background()
	r := NewHealthReporter().
		AddLivenessChecks(deadHealth("d")).
		AddReadinessChecks(deadHealth("d")).
		AddStatusChecks(deadHealth("d"))

	assert.Equal(t, StatusDead, r.LivenessStatus(ctx))
	assert.Equal(t, StatusUnReady, r.ReadinessStatus(ctx))
	assert.Equal(t, StatusDead, r.HealthReport(ctx).Alive)
}

func TestAliveAndReadyStatus(t *testing.T) {
	allOK := []HealthStatus{
		{Alive: StatusAlive, Ready: StatusReady},
		{Alive: StatusAlive, Ready: StatusReady},
	}
	assert.Equal(t, StatusAlive, AliveStatus(allOK...))
	assert.Equal(t, StatusReady, ReadyStatus(allOK...))

	mixed := []HealthStatus{
		{Alive: StatusAlive, Ready: StatusReady},
		{Alive: StatusDead, Ready: StatusUnReady},
	}
	assert.Equal(t, StatusDead, AliveStatus(mixed...))
	assert.Equal(t, StatusUnReady, ReadyStatus(mixed...))
}
