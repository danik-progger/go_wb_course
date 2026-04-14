package db

import (
	"sync"

	"urlshortener/domain"
)

type URLMemoryStore struct {
	mu   sync.RWMutex
	data map[string]domain.ShortenedURL
}

func NewURLMemoryStore() *URLMemoryStore {
	return &URLMemoryStore{
		data: make(map[string]domain.ShortenedURL),
	}
}

func (s *URLMemoryStore) Create(url domain.ShortenedURL) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[url.ID] = url
	return nil
}

func (s *URLMemoryStore) Get(id string) (domain.ShortenedURL, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.data[id]
	return url, ok
}

func (s *URLMemoryStore) Exists(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[id]
	return ok
}
