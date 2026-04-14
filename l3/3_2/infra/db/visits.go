package db

import (
	"sync"

	"urlshortener/domain"
)

type VisitMemoryStore struct {
	mu   sync.RWMutex
	data map[string][]domain.Visit
}

func NewVisitMemoryStore() *VisitMemoryStore {
	return &VisitMemoryStore{
		data: make(map[string][]domain.Visit),
	}
}

func (s *VisitMemoryStore) Record(visit domain.Visit) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[visit.ShortURLID] = append(s.data[visit.ShortURLID], visit)
}

func (s *VisitMemoryStore) GetByURL(id string) []domain.Visit {
	s.mu.RLock()
	defer s.mu.RUnlock()
	visits := s.data[id]
	if visits == nil {
		return []domain.Visit{}
	}
	result := make([]domain.Visit, len(visits))
	copy(result, visits)
	return result
}
