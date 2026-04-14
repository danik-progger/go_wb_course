package service

import (
	"fmt"
	"math/rand"
	"time"

	"urlshortener/app/repos"
	"urlshortener/domain"
)

type URLService struct {
	urls   repos.URLRepo
	visits repos.VisitRepo
	cache  repos.CacheRepo
}

func NewURLService(urls repos.URLRepo, visits repos.VisitRepo, cache repos.CacheRepo) *URLService {
	return &URLService{
		urls:   urls,
		visits: visits,
		cache:  cache,
	}
}

func (s *URLService) CreateShortURL(originalURL, customKey string) (string, error) {
	var shortKey string

	if customKey != "" {
		if s.urls.Exists(customKey) {
			return "", fmt.Errorf("custom short name is already taken")
		}
		shortKey = customKey
	} else {
		shortKey = generateShortKey()
		for s.urls.Exists(shortKey) {
			shortKey = generateShortKey()
		}
	}

	url := domain.ShortenedURL{
		ID:          shortKey,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
	}

	if err := s.urls.Create(url); err != nil {
		return "", fmt.Errorf("failed to create short URL: %w", err)
	}

	return shortKey, nil
}

func (s *URLService) Resolve(shortKey string) (string, error) {
	if originalURL, found := s.cache.Get(shortKey); found {
		return originalURL, nil
	}

	url, found := s.urls.Get(shortKey)
	if !found {
		return "", fmt.Errorf("shortened key not found")
	}

	s.cache.Set(shortKey, url.OriginalURL)
	return url.OriginalURL, nil
}

func (s *URLService) GetURL(shortKey string) (domain.ShortenedURL, bool) {
	return s.urls.Get(shortKey)
}

func (s *URLService) RecordVisit(shortKey, userAgent string) {
	s.visits.Record(domain.Visit{
		ShortURLID: shortKey,
		Timestamp:  time.Now(),
		UserAgent:  userAgent,
	})
}

func (s *URLService) GetVisits(shortKey string) []domain.Visit {
	return s.visits.GetByURL(shortKey)
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}
