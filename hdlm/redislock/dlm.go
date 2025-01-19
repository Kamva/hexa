package redislock

import (
	"context"
	"errors"
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"time"

	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
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
	ttl    time.Duration

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
	m.Expiry = time.Now().Add(m.ttl)

	_, err := m.client.Obtain(c, m.ID, m.ttl, nil)

	if errors.Is(err, redislock.ErrNotObtained) {
		err = hexa.ErrLockAlreadyAcquired
	}

	return tracer.Trace(err)
}

func (m *mutex) Unlock(c context.Context) error {
	return nil
}

var _ hexa.DLM = &dlm{}
var _ hexa.Mutex = &mutex{}
