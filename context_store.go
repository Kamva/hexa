package hexa

import "sync"

// Store is actually a concurrency-safe map.
type Store interface {
	Get(key string) interface{}
	Set(key string, val interface{})
}

type storeImpl struct {
	lock sync.RWMutex
	m    Map
}

func (s *storeImpl) Get(key string) interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.m[key]
}

func (s *storeImpl) Set(key string, val interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.m == nil {
		s.m = make(Map)
	}
	s.m[key] = val
}

func newStore() Store {
	return &storeImpl{}
}

var _ Store = &storeImpl{}
