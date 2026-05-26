package redislock

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/kamva/hexa"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newDlm builds a DLM whose redis client is not connected to anything; the
// constructor and mutex-building paths don't touch the network.
func newDlm(t *testing.T, owner string, ttl time.Duration) hexa.DLM {
	t.Helper()
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	t.Cleanup(func() { _ = client.Close() })

	d, err := NewDlm(DlmOptions{
		Client:          client,
		DefaultOwner:    owner,
		DefaultTTL:      ttl,
		WaitingInterval: 100 * time.Millisecond,
	})
	require.NoError(t, err)
	return d
}

func TestNewMutex_UsesDefaults(t *testing.T) {
	d := newDlm(t, "machine-1", time.Minute)
	m := d.NewMutex("key").(*mutex)

	assert.Equal(t, "key", m.ID)
	assert.Equal(t, "machine-1", m.Owner)
	assert.Equal(t, time.Minute, m.ttl)
}

func TestNewMutexWithTTL(t *testing.T) {
	d := newDlm(t, "machine-1", time.Minute)
	m := d.NewMutexWithTTL("key", 30*time.Second).(*mutex)

	assert.Equal(t, 30*time.Second, m.ttl)
	assert.Equal(t, "machine-1", m.Owner) // owner falls back to the dlm default
}

func TestNewMutexWithOptions_OwnerDefaulting(t *testing.T) {
	d := newDlm(t, "default-owner", time.Minute)

	withDefault := d.NewMutexWithOptions(hexa.MutexOptions{Key: "k", TTL: time.Second}).(*mutex)
	assert.Equal(t, "default-owner", withDefault.Owner)

	withExplicit := d.NewMutexWithOptions(hexa.MutexOptions{Key: "k2", Owner: "other", TTL: time.Second}).(*mutex)
	assert.Equal(t, "other", withExplicit.Owner)
}

// TestNewDlm_InitializesHealth is a regression test: the embedded Health used
// to be nil, so any health call on the DLM panicked with a nil-interface
// dispatch. It must now be a usable, non-nil Health.
func TestNewDlm_InitializesHealth(t *testing.T) {
	d := newDlm(t, "owner", time.Minute)

	h, ok := d.(hexa.Health)
	require.True(t, ok)
	assert.NotPanics(t, func() { _ = h.HealthIdentifier() })
	assert.Equal(t, "distributed_locks", h.HealthIdentifier())
}

// --- integration (needs a live redis) ---

func redisAddr(t *testing.T) string {
	addr := os.Getenv("HEXA_TEST_REDIS_ADDR")
	if addr == "" {
		t.Skip("skipping redislock integration test; set HEXA_TEST_REDIS_ADDR to run it")
	}
	return addr
}

func TestMutex_LockUnlock_Integration(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: redisAddr(t)})
	defer client.Close()

	d, err := NewDlm(DlmOptions{
		Client:          client,
		DefaultOwner:    "m1",
		DefaultTTL:      5 * time.Second,
		WaitingInterval: 100 * time.Millisecond,
	})
	require.NoError(t, err)

	ctx := context.Background()
	// Unique per run so the test doesn't collide with a concurrent run or a
	// stale key on a shared Redis instance.
	key := fmt.Sprintf("hexa-redislock-itest-%s-%d", t.Name(), time.Now().UnixNano())
	m1 := d.NewMutex(key)
	m2 := d.NewMutexWithOptions(hexa.MutexOptions{Key: key, Owner: "m2", TTL: 5 * time.Second})

	require.NoError(t, m1.Lock(ctx))

	// A different owner cannot acquire the held lock.
	assert.ErrorIs(t, m2.TryLock(ctx), hexa.ErrLockAlreadyAcquired)

	// Unlock releases it (the behavior that used to be a no-op).
	require.NoError(t, m1.Unlock(ctx))
	require.NoError(t, m2.TryLock(ctx))
	require.NoError(t, m2.Unlock(ctx))
}
