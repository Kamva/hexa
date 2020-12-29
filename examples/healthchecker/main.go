package main

import (
	"context"
	"syscall"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
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
		Tags:  map[string]string{"I'm": "ok :)"},
		Alive: h.LivenessStatus(ctx),
		Ready: h.ReadinessStatus(ctx),
	}
}

func main() {
	l := hlog.NewPrinterDriver(hlog.DebugLevel)
	checker := hexa.NewHealthChecker(l, "localhost:7676")

	r := hexa.NewHealthReporter().AddToChecks(&HealthExample{})

	gutil.PanicErr(checker.StartServer(r))
	gutil.WaitForSignals(syscall.SIGINT, syscall.SIGTERM)
}

var _ hexa.Health = &HealthExample{}
