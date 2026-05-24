package probe

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kamva/hexa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeHealth struct {
	alive hexa.LivenessStatus
	ready hexa.ReadinessStatus
}

func (f fakeHealth) HealthIdentifier() string                             { return "fake" }
func (f fakeHealth) LivenessStatus(context.Context) hexa.LivenessStatus   { return f.alive }
func (f fakeHealth) ReadinessStatus(context.Context) hexa.ReadinessStatus { return f.ready }
func (f fakeHealth) HealthStatus(context.Context) hexa.HealthStatus {
	return hexa.HealthStatus{Id: "fake", Alive: f.alive, Ready: f.ready}
}

// newProbe wires a probe server's mux behind an httptest server.
func newProbe(t *testing.T, h hexa.Health) *httptest.Server {
	mux := http.NewServeMux()
	ps := NewServer(&http.Server{}, mux)
	RegisterHealthHandlers(ps, hexa.NewHealthReporter().AddToChecks(h))
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

func TestHealthHandlers_Healthy(t *testing.T) {
	ts := newProbe(t, fakeHealth{alive: hexa.StatusAlive, ready: hexa.StatusReady})

	for _, p := range []string{"/live", "/ready", "/status"} {
		resp, err := http.Get(ts.URL + p)
		require.NoError(t, err, p)
		assert.Equal(t, http.StatusOK, resp.StatusCode, p)
		_ = resp.Body.Close()
	}
}

func TestLivenessHandler_Dead(t *testing.T) {
	ts := newProbe(t, fakeHealth{alive: hexa.StatusDead, ready: hexa.StatusUnReady})

	resp, err := http.Get(ts.URL + "/live")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, string(hexa.StatusDead), resp.Header.Get("liveness_status"))
}

func TestReadinessHandler_Unready(t *testing.T) {
	ts := newProbe(t, fakeHealth{alive: hexa.StatusAlive, ready: hexa.StatusUnReady})

	resp, err := http.Get(ts.URL + "/ready")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

func TestStatusHandler_ReturnsReport(t *testing.T) {
	ts := newProbe(t, fakeHealth{alive: hexa.StatusAlive, ready: hexa.StatusReady})

	resp, err := http.Get(ts.URL + "/status")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var body struct {
		Code string            `json:"code"`
		Data hexa.HealthReport `json:"data"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, "app.status", body.Code)
	assert.Equal(t, hexa.StatusAlive, body.Data.Alive)
	assert.Equal(t, hexa.StatusReady, body.Data.Ready)
	require.Len(t, body.Data.Statuses, 1)
	assert.Equal(t, "fake", body.Data.Statuses[0].Id)
}

func TestDocsHandler_ListsRegisteredHandlers(t *testing.T) {
	mux := http.NewServeMux()
	ps := NewServer(&http.Server{}, mux)
	ps.Register("custom", "/custom", func(http.ResponseWriter, *http.Request) {}, "a custom handler")
	ts := httptest.NewServer(mux)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/")
	require.NoError(t, err)
	defer resp.Body.Close()

	var ds []HandlerDescriptor
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&ds))

	names := make([]string, len(ds))
	for i, d := range ds {
		names[i] = d.Name
	}
	assert.Contains(t, names, "docs")
	assert.Contains(t, names, "custom")
}

func TestServer_RunAndShutdown(t *testing.T) {
	ps := NewServer(&http.Server{Addr: "127.0.0.1:0"}, http.NewServeMux())

	done, err := ps.Run()
	require.NoError(t, err)

	require.NoError(t, ps.Shutdown(context.Background()))

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("done channel not closed after shutdown")
	}
}
