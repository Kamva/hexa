package hexa

import "context"

type ReadinessStatus string

const (
	StatusReady   ReadinessStatus = "READY"
	StatusUnReady ReadinessStatus = "UNREADY"
)

type LivenessStatus string

const (
	StatusAlive = "alive"
	StatusDead  = "dead"
)

type Health interface {
	HealthIdentifier() string
	LivenessStatus(ctx context.Context) LivenessStatus
	ReadinessStatus(ctx context.Context) ReadinessStatus
}

func LivenessCheck(l ...Health) map[string]LivenessStatus {
	r := make(map[string]LivenessStatus, len(l))
	for _, health := range l {
		r[health.HealthIdentifier()] = health.LivenessStatus(context.Background())
	}
	return r
}

func ReadinessCheck(l ...Health) map[string]ReadinessStatus {
	// TODO: check using go routine
	r := make(map[string]ReadinessStatus, len(l))
	for _, health := range l {
		r[health.HealthIdentifier()] = health.ReadinessStatus(context.Background())
	}
	return r
}

type HealthCheckResult struct {
	Live  LivenessStatus  `json:"live"`
	Ready ReadinessStatus `json:"ready"`
}

func HealthCheck(l ...Health) map[string]HealthCheckResult {
	// TODO: check using go routines
	r := make(map[string]HealthCheckResult, len(l))
	for _, health := range l {
		r[health.HealthIdentifier()] = HealthCheckResult{
			Live:  health.LivenessStatus(context.Background()),
			Ready: health.ReadinessStatus(context.Background()),
		}
	}

	return r
}
