package inmemorystorage

import (
	"context"
	"sync"
)

type set map[string]struct{}

type Storage struct {
	t   map[string]set
	tMu *sync.RWMutex
}

func New() *Storage {
	return &Storage{
		t:   map[string]set{},
		tMu: &sync.RWMutex{},
	}
}

func (s *Storage) Fetch(_ context.Context, key string) ([]string, error) {
	s.tMu.RLock()
	defer s.tMu.RUnlock()

	v, ok := s.t[key]
	if !ok {
		return nil, nil
	}

	result := make([]string, 0, len(v))

	for k := range v {
		result = append(result, k)
	}

	return result, nil
}

func (s *Storage) Insert(_ context.Context, key string, value string) error {
	s.tMu.Lock()
	defer s.tMu.Unlock()

	v, ok := s.t[key]
	if !ok {
		s.t[key] = set{value: {}}

		return nil
	}

	v[value] = struct{}{}
	s.t[key] = v

	return nil
}
