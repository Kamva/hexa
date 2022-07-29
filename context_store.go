package hexa

import (
	"sync"
)

// Store is actually a concurrency-safe map.
type Store interface {
	Get(key string) any
	Set(key string, val any)
	SetIfNotExist(key string, val func() any) any
}

type atomicStore struct {
	lock sync.RWMutex
	m    map[string]any
}

func (s *atomicStore) Get(key string) any {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.m[key]
}

func (s *atomicStore) Set(key string, val any) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.m == nil {
		s.m = make(map[string]any)
	}
	s.m[key] = val
}

func (s *atomicStore) SetIfNotExist(key string, vp func() any) any {
	s.lock.RLock()
	val := s.m[key]
	if val != nil {
		s.lock.RUnlock()
		return val
	}

	s.lock.RUnlock()
	s.lock.Lock()
	defer s.lock.Unlock()

	val = s.m[key] // check if exists again, maybe when we were changing the locks, someone set the value.
	if val != nil {
		return val
	}

	if s.m == nil {
		s.m = make(map[string]any)
	}

	val = vp()
	s.m[key] = val
	return val
}

func newStore() Store {
	return &atomicStore{}
}

var _ Store = &atomicStore{}
