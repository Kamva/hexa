package sr

import (
	"context"
	"errors"
	"testing"

	"github.com/kamva/hexa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// recorder captures lifecycle call order across the fake services.
type recorder struct {
	booted   []string
	shutdown []string
}

// svc is a fake service that can be Bootable and Shutdownable.
type svc struct {
	name    string
	rec     *recorder
	bootErr error
}

func (s *svc) Boot() error {
	s.rec.booted = append(s.rec.booted, s.name)
	return s.bootErr
}

func (s *svc) Shutdown(context.Context) error {
	s.rec.shutdown = append(s.rec.shutdown, s.name)
	return nil
}

// healthSvc additionally implements hexa.Health.
type healthSvc struct {
	id string
}

func (h *healthSvc) HealthIdentifier() string                          { return h.id }
func (h *healthSvc) LivenessStatus(context.Context) hexa.LivenessStatus { return hexa.StatusAlive }
func (h *healthSvc) ReadinessStatus(context.Context) hexa.ReadinessStatus {
	return hexa.StatusReady
}
func (h *healthSvc) HealthStatus(context.Context) hexa.HealthStatus {
	return hexa.HealthStatus{Id: h.id}
}

func names(ds []*hexa.Descriptor) []string {
	out := make([]string, len(ds))
	for i, d := range ds {
		out[i] = d.Name
	}
	return out
}

func TestRegisterByDescriptor_OrdersByPriority(t *testing.T) {
	r := New()
	r.RegisterByDescriptor(&hexa.Descriptor{Name: "a", Instance: &struct{}{}, Priority: 3})
	r.RegisterByDescriptor(&hexa.Descriptor{Name: "b", Instance: &struct{}{}, Priority: 1})
	r.RegisterByDescriptor(&hexa.Descriptor{Name: "c", Instance: &struct{}{}, Priority: 2})

	assert.Equal(t, []string{"b", "c", "a"}, names(r.Descriptors()))
}

func TestRegister_DefaultPriorityIsNextValue(t *testing.T) {
	r := New()
	r.Register("a", &struct{}{})
	r.Register("b", &struct{}{})

	assert.Equal(t, []string{"a", "b"}, names(r.Descriptors()))
	assert.Equal(t, 1, r.Descriptor("a").Priority)
	assert.Equal(t, 2, r.Descriptor("b").Priority)
}

func TestRegisterByInstance_UsesTypeName(t *testing.T) {
	r := New()
	r.RegisterByInstance(&svc{name: "ignored", rec: &recorder{}})

	// reflect type name of *svc is "svc".
	require.NotNil(t, r.Descriptor("svc"))
	assert.NotNil(t, r.Service("svc"))
}

func TestServiceAndDescriptor_MissingReturnsNil(t *testing.T) {
	r := New()
	assert.Nil(t, r.Service("missing"))
	assert.Nil(t, r.Descriptor("missing"))
}

func TestBoot_RunsBootablesInPriorityOrder(t *testing.T) {
	rec := &recorder{}
	r := New()
	r.RegisterByDescriptor(&hexa.Descriptor{Name: "a", Instance: &svc{name: "a", rec: rec}, Priority: 2})
	r.RegisterByDescriptor(&hexa.Descriptor{Name: "b", Instance: &svc{name: "b", rec: rec}, Priority: 1})

	require.NoError(t, r.Boot())
	assert.Equal(t, []string{"b", "a"}, rec.booted)

	// Boot is idempotent: a second call does not re-boot.
	require.NoError(t, r.Boot())
	assert.Equal(t, []string{"b", "a"}, rec.booted)
}

func TestBoot_ReturnsBootError(t *testing.T) {
	rec := &recorder{}
	r := New()
	bootErr := errors.New("boom")
	r.Register("a", &svc{name: "a", rec: rec, bootErr: bootErr})

	err := r.Boot()
	require.Error(t, err)
	assert.ErrorIs(t, err, bootErr)
}

func TestShutdown_RunsShutdownablesInReverseOrder(t *testing.T) {
	rec := &recorder{}
	r := New()
	r.RegisterByDescriptor(&hexa.Descriptor{Name: "a", Instance: &svc{name: "a", rec: rec}, Priority: 1})
	r.RegisterByDescriptor(&hexa.Descriptor{Name: "b", Instance: &svc{name: "b", rec: rec}, Priority: 2})

	require.NoError(t, r.Shutdown(context.Background()))
	// Highest priority shuts down first.
	assert.Equal(t, []string{"b", "a"}, rec.shutdown)

	// ShutdownCh is closed once shutdown completed.
	select {
	case <-r.ShutdownCh():
	default:
		t.Fatal("expected ShutdownCh to be closed after Shutdown")
	}
}

func TestRegisterByDescriptor_DetectsHealthFromInstance(t *testing.T) {
	r := New()
	r.RegisterByDescriptor(&hexa.Descriptor{Name: "h", Instance: &healthSvc{id: "h"}})

	d := r.Descriptor("h")
	require.NotNil(t, d)
	require.NotNil(t, d.Health)
	assert.Equal(t, "h", d.Health.HealthIdentifier())
}

func TestRegisterByDescriptor_KeepsExplicitHealth(t *testing.T) {
	explicit := &healthSvc{id: "explicit"}
	r := New()
	// Instance also implements Health, but an explicit Health must win.
	r.RegisterByDescriptor(&hexa.Descriptor{
		Name:     "h",
		Instance: &healthSvc{id: "instance"},
		Health:   explicit,
	})

	assert.Equal(t, "explicit", r.Descriptor("h").Health.HealthIdentifier())
}

func TestMultiSearchRegistry_SearchesAllRegistries(t *testing.T) {
	primary := New()

	r1 := New()
	r1.Register("one", &struct{}{})
	r2 := New()
	r2.Register("two", &struct{}{})

	multi := NewMultiSearchRegistry(primary, r1, r2)

	// Descriptor/Service search across every registry, including the first.
	require.NotNil(t, multi.Descriptor("one"))
	require.NotNil(t, multi.Descriptor("two"))
	assert.NotNil(t, multi.Service("one"))
	assert.Nil(t, multi.Service("missing"))

	// Descriptors must include both registries' descriptors (the first
	// registry used to be dropped by an off-by-one).
	got := names(multi.Descriptors())
	assert.ElementsMatch(t, []string{"one", "two"}, got)
}
