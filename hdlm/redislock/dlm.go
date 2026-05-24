package redislock

import (
	"context"
	"errors"
	"time"

	"github.com/bsm/redislock"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"github.com/redis/go-redis/v9"
)

type DlmOptions struct {
	Client          *redis.Client
	WaitingInterval time.Duration
	DefaultTTL      time.Duration
	DefaultOwner    string
}

// dlm implements the Hexa DLM.
type dlm struct {
	hexa.Health
	client   *redislock.Client
	owner    string
	ttl      time.Duration
	interval time.Duration
}

func NewDlm(o DlmOptions) (hexa.DLM, error) {
	dlm := &dlm{
		Health: hexa.NewPingHealth(hlog.GlobalLogger(), "distributed_locks", func(ctx context.Context) error {
			return o.Client.Ping(ctx).Err()
		}, nil),

		client:   redislock.New(o.Client),
		ttl:      o.DefaultTTL,
		owner:    o.DefaultOwner,
		interval: o.WaitingInterval,
	}

	return dlm, nil
}

func (m *dlm) NewMutex(Key string) hexa.Mutex {
	return m.NewMutexWithOptions(hexa.MutexOptions{
		Key:   Key,
		Owner: m.owner,
		TTL:   m.ttl,
	})
}

func (m *dlm) NewMutexWithTTL(Key string, ttl time.Duration) hexa.Mutex {
	return m.NewMutexWithOptions(hexa.MutexOptions{Key: Key, TTL: ttl})
}

func (m *dlm) NewMutexWithOptions(o hexa.MutexOptions) hexa.Mutex {
	if o.Owner == "" {
		o.Owner = m.owner
	}

	return &mutex{
		client:   m.client,
		ttl:      o.TTL,
		interval: m.interval,

		ID:    o.Key,
		Owner: o.Owner,
	}
}

type mutex struct {
	client *redislock.Client
	// lock holds the currently acquired redis lock. It is set on a
	// successful (Try)Lock and cleared on Unlock, and carries the token
	// required to release or refresh the lock.
	lock *redislock.Lock
	ttl  time.Duration

	// interval is waiting time before try to lock again if
	// lock already held by another mutex.
	interval time.Duration

	ID    string `json:"key" bson:"_id"`
	Owner string `json:"owner" bson:"owner"`
	// Expiry begins when we lock the mutex.
	Expiry time.Time `json:"expiry" bson:"expiry"`
}

func (m *mutex) Lock(c context.Context) error {
	for {
		err := m.TryLock(c)

		if !errors.Is(err, hexa.ErrLockAlreadyAcquired) {
			return tracer.Trace(err)
		}

		select {
		case <-c.Done():
			return c.Err()
		case <-time.After(m.interval):
		}
	}
}

func (m *mutex) TryLock(c context.Context) error {
	// Use the owner as the lock token so mutexes that share an owner can
	// release/refresh each other's lock (an empty owner yields a random
	// token, giving each mutex its own independent lock).
	opts := &redislock.Options{Token: m.Owner}

	// If we already hold the lock, refresh it instead of failing so repeated
	// (Try)Lock calls extend the lease rather than reporting it as taken.
	if m.lock != nil {
		err := m.lock.Refresh(c, m.ttl, opts)
		if err == nil {
			m.Expiry = time.Now().Add(m.ttl)
			return nil
		}
		if !errors.Is(err, redislock.ErrNotObtained) {
			return tracer.Trace(err)
		}
		m.lock = nil // We lost the lock; fall through and try to obtain it again.
	}

	lock, err := m.client.Obtain(c, m.ID, m.ttl, opts)
	if errors.Is(err, redislock.ErrNotObtained) {
		return tracer.Trace(hexa.ErrLockAlreadyAcquired)
	}
	if err != nil {
		return tracer.Trace(err)
	}

	m.lock = lock
	m.Expiry = time.Now().Add(m.ttl)
	return nil
}

func (m *mutex) Unlock(c context.Context) error {
	if m.lock == nil {
		return nil
	}

	err := m.lock.Release(c)
	m.lock = nil
	if errors.Is(err, redislock.ErrLockNotHeld) {
		// Already released or expired: a no-op per the Mutex contract.
		return nil
	}
	return tracer.Trace(err)
}

var _ hexa.DLM = &dlm{}
var _ hexa.Mutex = &mutex{}
