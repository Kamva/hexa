package probe

import (
	"fmt"
	"net/http"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
)

const (
	livenessStatusKey  = "liveness_status"
	readinessStatusKey = "readiness_status"
)

type healthHandlers struct {
	r hexa.HealthReporter
}

func RegisterHealthHandlers(ps Server, r hexa.HealthReporter) {
	(&healthHandlers{r: r}).RegisterHandlers(ps)
}

func (h *healthHandlers) RegisterHandlers(ps Server) {
	ps.Register("live", "/live", h.livenessHandler, "reports app's liveness")
	ps.Register("ready", "/ready", h.readinessHandler, "reports app's readiness")
	ps.Register("status", "/status", h.statusHandler, "reports app's status")
}

func (h *healthHandlers) livenessHandler(w http.ResponseWriter, r *http.Request) {
	status := h.r.LivenessStatus(r.Context())
	w.Header().Set(livenessStatusKey, string(status))

	if status != hexa.StatusAlive {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *healthHandlers) readinessHandler(w http.ResponseWriter, r *http.Request) {
	status := h.r.ReadinessStatus(r.Context())
	w.Header().Set(readinessStatusKey, string(status))

	if status != hexa.StatusReady {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *healthHandlers) statusHandler(w http.ResponseWriter, r *http.Request) {
	report := h.r.HealthReport(r.Context())
	w.Header().Set(livenessStatusKey, string(report.Alive))
	w.Header().Set(readinessStatusKey, string(report.Ready))
	w.Header().Set("Content-Type", "application/json")

	resp := hexa.Map{
		"code": "app.status",
		"data": report,
	}

	b, err := gutil.Marshal(resp)
	if err != nil {
		hlog.Error("error on marshaling health report", hlog.ErrStack(tracer.Trace(err)))
		resp := fmt.Sprintf(`{"err" : "%s"}`, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(resp))
		return
	}

	w.Write(b)
}
