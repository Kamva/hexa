package hexa

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kamva/gutil"
	"github.com/kamva/tracer"
)

var ()

const (
	LivenessStatusKey  = "liveness_status"
	ReadinessStatusKey = "readiness_status"
)

type HealthChecker interface {
	HealthCheck(l ...Health) []HealthStatus
	StartHealthCheckServer(hp HealthProbe) error
	StopHealthCheckServer() error
}

type healthChecker struct {
	l      Logger
	server *http.Server
	addr   string
}

type HealthCheckerOptions struct {
	Logger  Logger
	Address string
}

func NewHealthChecker(o HealthCheckerOptions) HealthChecker {
	return &healthChecker{
		l:    o.Logger,
		addr: o.Address,
	}
}

func (h *healthChecker) HealthCheck(l ...Health) []HealthStatus {
	// TODO: check using go routines
	r := make([]HealthStatus, len(l))
	for i, health := range l {
		r[i] = HealthStatus{
			Id:    health.HealthIdentifier(),
			Live:  health.LivenessStatus(context.Background()),
			Ready: health.ReadinessStatus(context.Background()),
		}
	}

	return r
}

func (h *healthChecker) StartHealthCheckServer(hp HealthProbe) error {
	if h.server != nil {
		if err := h.server.Shutdown(context.Background()); err != nil && err != http.ErrServerClosed {
			return tracer.Trace(err)
		}
	}

	mux := http.NewServeMux()
	h.server = &http.Server{Addr: h.addr, Handler: mux}

	mux.HandleFunc("/live", h.livenessHandler(hp))
	mux.HandleFunc("/ready", h.readinessHandler(hp))
	mux.HandleFunc("/status", h.statusHandler(hp))

	h.l.Info("start serving health check requests", StringField("address", h.addr))
	go func() {
		h.server.ListenAndServe()
	}()

	return nil
}

func (h *healthChecker) StopHealthCheckServer() error {
	return h.server.Close()
}

//--------------------------------
// HTTP Health Check Handlers
//--------------------------------

func (h *healthChecker) livenessHandler(hp HealthProbe) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		status := hp.LivenessStatus(r.Context())
		w.Header().Set(LivenessStatusKey, string(status))

		if status != StatusAlive {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *healthChecker) readinessHandler(hp HealthProbe) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		status := hp.ReadinessStatus(r.Context())
		w.Header().Set(ReadinessStatusKey, string(status))

		if status != StatusReady {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *healthChecker) statusHandler(hp HealthProbe) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		report := hp.HealthReport(r.Context())
		w.Header().Set(LivenessStatusKey, string(report.Live))
		w.Header().Set(ReadinessStatusKey, string(report.Ready))
		//fmt.Fprint(w, gutil.UnmarshalStruct())
		//return w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		resp := Map{
			"code":   "app.status",
			"data": report,
		}

		b, err := gutil.Marshal(resp)
		if err != nil {
			h.l.Error("error on marshaling health report", ErrField(err), ErrStackField(tracer.Trace(err)))
			resp := fmt.Sprintf(`{"err" : "%s"}`, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(resp))
			return
		}

		w.Write(b)
	}
}

var _ HealthChecker = &healthChecker{}
