package hexa

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kamva/gutil"
	"github.com/kamva/tracer"
)

const (
	LivenessStatusKey  = "liveness_status"
	ReadinessStatusKey = "readiness_status"
)

type HealthChecker interface {
	Runnable
	Shutdownable
}

type healthChecker struct {
	l      Logger
	server *http.Server
	addr   string
	r      HealthReporter
}

func NewHealthChecker(l Logger, addr string, r HealthReporter) HealthChecker {
	return &healthChecker{
		l:    l,
		addr: addr,
		r:    r,
	}
}

func (h *healthChecker) Run() error {
	if h.server != nil {
		if err := h.server.Shutdown(context.Background()); err != nil && err != http.ErrServerClosed {
			return tracer.Trace(err)
		}
	}

	mux := http.NewServeMux()
	h.server = &http.Server{Addr: h.addr, Handler: mux}

	mux.HandleFunc("/live", h.livenessHandler())
	mux.HandleFunc("/ready", h.readinessHandler())
	mux.HandleFunc("/status", h.statusHandler())

	h.l.Info("start serving health check requests", StringField("address", h.addr))
	go func() {
		err := h.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			h.l.Error("error on health check server", ErrStackField(tracer.Trace(err)), ErrField(err))
		}
	}()

	return nil
}

func (h *healthChecker) Shutdown(c context.Context) error {
	return h.server.Shutdown(c)
}

//--------------------------------
// HTTP Health Check Handlers
//--------------------------------

func (h *healthChecker) livenessHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		status := h.r.LivenessStatus(r.Context())
		w.Header().Set(LivenessStatusKey, string(status))

		if status != StatusAlive {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *healthChecker) readinessHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		status := h.r.ReadinessStatus(r.Context())
		w.Header().Set(ReadinessStatusKey, string(status))

		if status != StatusReady {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *healthChecker) statusHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		report := h.r.HealthReport(r.Context())
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
