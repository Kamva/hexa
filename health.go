package hexa

import (
	"context"
)

type ReadinessStatus string

const (
	StatusReady   ReadinessStatus = "READY"
	StatusUnReady ReadinessStatus = "UNREADY"
)

type LivenessStatus string

const (
	StatusAlive LivenessStatus = "ALIVE"
	StatusDead  LivenessStatus = "DEAD"
)

type (
	LivenessResult struct {
		Id     string `json:"id"`
		Status LivenessStatus
	}

	ReadinessResult struct {
		Id     string `json:"id"`
		Status ReadinessStatus
	}

	HealthStatus struct {
		Id    string            `json:"id"`
		Tags  map[string]string `json:"tags,omitempty"`
		Alive LivenessStatus    `json:"alive"`
		Ready ReadinessStatus   `json:"ready"`
	}
)

type HealthReport struct {
	HealthStatus
	Statuses []HealthStatus `json:"statuses"`
}

type Health interface {
	HealthIdentifier() string
	LivenessStatus(ctx context.Context) LivenessStatus
	ReadinessStatus(ctx context.Context) ReadinessStatus
	HealthStatus(ctx context.Context) HealthStatus
}

type HealthReporter interface {
	LivenessStatus(ctx context.Context) LivenessStatus
	ReadinessStatus(ctx context.Context) ReadinessStatus
	HealthReport(ctx context.Context) HealthReport
}

func HealthCheck(l ...Health) []HealthStatus {
	// TODO: check using go routines
	r := make([]HealthStatus, len(l))
	for i, health := range l {
		r[i] = HealthStatus{
			Id:    health.HealthIdentifier(),
			Alive: health.LivenessStatus(context.Background()),
			Ready: health.ReadinessStatus(context.Background()),
		}
	}

	return r
}

func AllAliveStatus(l ...HealthStatus) bool {
	for _, s := range l {
		if s.Alive != StatusAlive {
			return false
		}
	}
	return true
}

func AllReadyStatus(l ...HealthStatus) bool {
	for _, s := range l {
		if s.Ready != StatusReady {
			return false
		}
	}
	return true
}
