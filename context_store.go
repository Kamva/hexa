package hexa

import (
	"sync"
)

// Store is actually a concurrency-safe map.
type Store interface {
	Get(key string) interface{}
	Set(key string, val interface{})
	SetIfNotExist(key string, val func() interface{}) interface{}
}

type atomicStore struct {
	lock sync.RWMutex
	m    map[string]interface{}
}

func (s *atomicStore) Get(key string) interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.m[key]
}

func (s *atomicStore) Set(key string, val interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.m[key] = val
}

func (s *atomicStore) SetIfNotExist(key string, vp func() interface{}) interface{} {
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

	val = vp()
	s.m[key] = val
	return val
}

func newStore() Store {
	return &atomicStore{}
}

var _ Store = &atomicStore{}
