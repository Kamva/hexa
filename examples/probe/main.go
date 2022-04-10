package main

import (
	"context"
	"net/http"
	"syscall"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/probe"
)

type HealthExample struct {
}

func (h *HealthExample) HealthIdentifier() string {
	return "health_example"
}

func (h *HealthExample) LivenessStatus(ctx context.Context) hexa.LivenessStatus {
	return hexa.StatusAlive
}

func (h *HealthExample) ReadinessStatus(ctx context.Context) hexa.ReadinessStatus {
	return hexa.StatusReady
}

func (h *HealthExample) HealthStatus(ctx context.Context) hexa.HealthStatus {
	return hexa.HealthStatus{
		Id:    h.HealthIdentifier(),
		Tags:  map[string]string{"I'm": "good :)"},
		Alive: h.LivenessStatus(ctx),
		Ready: h.ReadinessStatus(ctx),
	}
}

func main() {
	r := hexa.NewHealthReporter().AddToChecks(&HealthExample{})

	s := probe.NewServer(&http.Server{Addr: "localhost:7676"}, http.NewServeMux())
	probe.RegisterHealthHandlers(s, r)
	probe.RegisterPprofHandlers(s)

	gutil.PanicErr(s.Run())
	gutil.WaitForSignals(syscall.SIGINT, syscall.SIGTERM)
}

var _ hexa.Health = &HealthExample{}
