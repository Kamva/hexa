package hexa

import (
	"sync"
)

// Store is actually a concurrency-safe map.
type Store interface {
	Get(key string) interface{}
	Set(key string, val interface{})
	// Atomic runs a function in atomic mode.
	// You can not call to the Atomic function inside an atmoic function. it will panic.
	Atomic(func(s Store))
}

type atomicStore struct {
	lock sync.RWMutex
	m    *mapStore
}

func (s *atomicStore) Get(key string) interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.m.Get(key)
}

func (s *atomicStore) Set(key string, val interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.m.Set(key, val)
}

func (s *atomicStore) Atomic(f func(s Store)) {
	s.lock.Lock()
	defer s.lock.Unlock()
	f(s.m)
}

func newStore() Store {
	return &atomicStore{m: &mapStore{}}
}

// mapStore implements the store but doesn't support Atomic mode.
type mapStore struct {
	m map[string]interface{}
}

func (s *mapStore) Get(key string) interface{} {
	return s.m[key]
}

func (s *mapStore) Set(key string, val interface{}) {
	if s.m == nil {
		s.m = make(Map)
	}
	s.m[key] = val
}

func (s *mapStore) Atomic(f func(s Store)) {
	panic("map store doesn't support atomic function")
}

var _ Store = &atomicStore{}
var _ Store = &mapStore{}
