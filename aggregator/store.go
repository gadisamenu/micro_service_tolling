package main

import (
	"fmt"

	"github.com/gadisamenu/tolling/types"
)

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (s *MemoryStore) Insert(distance types.Distance) error {
	s.data[distance.ObuId] += distance.Value
	return nil
}

func (s *MemoryStore) Get(obuId int) (float64, error) {
	value, ok := s.data[obuId]
	if !ok {
		return 0.0, fmt.Errorf("distance not found with obu id %d", obuId)
	}
	return value, nil
}
