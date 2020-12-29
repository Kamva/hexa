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
	StartServer(r HealthReporter) error
	StopServer() error
}

type healthChecker struct {
	l      Logger
	server *http.Server
	addr   string
}

func NewHealthChecker(l Logger, addr string) HealthChecker {
	return &healthChecker{
		l:    l,
		addr: addr,
	}
}

func (h *healthChecker) StartServer(r HealthReporter) error {
	if h.server != nil {
		if err := h.server.Shutdown(context.Background()); err != nil && err != http.ErrServerClosed {
			return tracer.Trace(err)
		}
	}

	mux := http.NewServeMux()
	h.server = &http.Server{Addr: h.addr, Handler: mux}

	mux.HandleFunc("/live", h.livenessHandler(r))
	mux.HandleFunc("/ready", h.readinessHandler(r))
	mux.HandleFunc("/status", h.statusHandler(r))

	h.l.Info("start serving health check requests", StringField("address", h.addr))
	go func() {
		err := h.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			h.l.Error("error on health check server", ErrStackField(tracer.Trace(err)), ErrField(err))
		}
	}()

	return nil
}

func (h *healthChecker) StopServer() error {
	return h.server.Close()
}

//--------------------------------
// HTTP Health Check Handlers
//--------------------------------

func (h *healthChecker) livenessHandler(hp HealthReporter) func(http.ResponseWriter, *http.Request) {
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

func (h *healthChecker) readinessHandler(hp HealthReporter) func(http.ResponseWriter, *http.Request) {
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

func (h *healthChecker) statusHandler(hp HealthReporter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		report := hp.HealthReport(r.Context())
		w.Header().Set(LivenessStatusKey, string(report.Alive))
		w.Header().Set(ReadinessStatusKey, string(report.Ready))
		//fmt.Fprint(w, gutil.UnmarshalStruct())
		//return w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		resp := Map{
			"code": "app.status",
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
