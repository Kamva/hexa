package hexa

import (
	"context"
	"errors"
	"time"
)

var ErrLockAlreadyAcquired = errors.New("lock already acquired")

type MutexOptions struct {
	Key string
	// If owner is empty, DLM uses default Owner.
	Owner string
	TTL   time.Duration
}

// DLM is distributed lock manager.
type DLM interface {
	// NewMutex returns new mutex with provided key
	// and default owner and ttl.
	NewMutex(Key string) Mutex

	// NewMutexWithTTL creates a new mutex with
	// the provided key and ttl.
	NewMutexWithTTL(key string, ttl time.Duration) Mutex

	// NewMutexWithOptions returns new mutex with
	// provided options.
	NewMutexWithOptions(o MutexOptions) Mutex
}

// Mutex can be used as a Distributed lock.
type Mutex interface {
	// Lock try to lock or wait for release and then lock.
	// We can have multiple behaviors when Lock invoke multiple times after locked one time:
	// 1. try to refresh the lock
	// 2. return error
	// 3. return nil and ignore next calls.
	// behavior in our implementation should be 1, means
	// you should try to refresh lock when user call
	// this method again.
	Lock(ctx context.Context) error

	// TryLock tries to lock or returns the ErrLockAlreadyAcquired
	// error if it acquired.
	// Please note different mutexes with same lock name and same
	// owner can lock and unlock each other.
	// expiry date should begin after call to this method, not at the
	// creation time of this mutex.
	// We can have multiple behaviors when TryLock invoke multiple times after locked one time:
	// 1. try to refresh the lock
	// 2. return error
	// 3. return nil and ignore next calls.
	// behavior in our implementation should be the option number 1, means
	// you should try to refresh lock when user calls to
	// this method again.
	TryLock(ctx context.Context) error

	// Unlock release the lock.
	// it should ignore if lock is already released
	// and do not return any error.
	Unlock(ctx context.Context) error
}
