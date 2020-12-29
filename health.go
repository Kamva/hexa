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
	Alive    LivenessStatus  `json:"alive"`
	Ready    ReadinessStatus `json:"ready"`
	Statuses []HealthStatus  `json:"statuses"`
}

type Health interface {
	HealthIdentifier() string
	LivenessStatus(ctx context.Context) LivenessStatus
	ReadinessStatus(ctx context.Context) ReadinessStatus
	HealthStatus(ctx context.Context) HealthStatus
}

type HealthReporter interface {
	AddLivenessChecks(l ...Health) HealthReporter
	AddReadinessChecks(l ...Health) HealthReporter
	AddStatusChecks(l ...Health) HealthReporter
	AddToChecks(l ...Health) HealthReporter
	LivenessStatus(ctx context.Context) LivenessStatus
	ReadinessStatus(ctx context.Context) ReadinessStatus
	HealthReport(ctx context.Context) HealthReport
}

type healthReporter struct {
	livenssCheck   []Health
	readinessCheck []Health
	statusCheck    []Health
}

func NewHealthReporter() HealthReporter {
	return &healthReporter{
		livenssCheck:   []Health{},
		readinessCheck: []Health{},
		statusCheck:    []Health{},
	}
}

func (h *healthReporter) AddLivenessChecks(l ...Health) HealthReporter {
	h.livenssCheck = append(h.livenssCheck, l...)
	return h
}

func (h *healthReporter) AddReadinessChecks(l ...Health) HealthReporter {
	h.readinessCheck = append(h.readinessCheck, l...)
	return h
}

func (h *healthReporter) AddStatusChecks(l ...Health) HealthReporter {
	h.statusCheck = append(h.statusCheck, l...)
	return h
}

func (h *healthReporter) AddToChecks(l ...Health) HealthReporter {
	return h.AddLivenessChecks(l...).AddReadinessChecks(l...).AddStatusChecks(l...)
}

func (h healthReporter) LivenessStatus(ctx context.Context) LivenessStatus {
	for _, health := range h.livenssCheck {
		if st := health.LivenessStatus(ctx); st != StatusAlive {
			return st
		}
	}
	return StatusAlive
}

func (h healthReporter) ReadinessStatus(ctx context.Context) ReadinessStatus {
	for _, health := range h.readinessCheck {
		if st := health.ReadinessStatus(ctx); st != StatusReady {
			return st
		}
	}
	return StatusReady
}

func (h healthReporter) HealthReport(ctx context.Context) HealthReport {
	l := HealthCheck(h.statusCheck...)
	return HealthReport{
		Alive:    AliveStatus(l...),
		Ready:    ReadyStatus(l...),
		Statuses: l,
	}
}

var _ HealthReporter = &healthReporter{}

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

func AliveStatus(l ...HealthStatus) LivenessStatus {
	for _, s := range l {
		if s.Alive != StatusAlive {
			return StatusDead
		}
	}
	return StatusAlive
}

func ReadyStatus(l ...HealthStatus) ReadinessStatus {
	for _, s := range l {
		if s.Ready != StatusReady {
			return StatusUnReady
		}
	}
	return StatusReady
}
