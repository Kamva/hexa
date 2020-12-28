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
		Live  LivenessStatus    `json:"live"`
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

type HealthProbe interface {
	LivenessStatus(ctx context.Context) LivenessStatus
	ReadinessStatus(ctx context.Context) ReadinessStatus
	HealthReport(ctx context.Context) HealthReport
}
