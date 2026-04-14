package service

import (
	"testing"
	"time"

	"urlshortener/domain"
)

type mockURLRepo struct {
	data map[string]domain.ShortenedURL
}

func newMockURLRepo() *mockURLRepo {
	return &mockURLRepo{data: make(map[string]domain.ShortenedURL)}
}

func (r *mockURLRepo) Create(url domain.ShortenedURL) error {
	r.data[url.ID] = url
	return nil
}

func (r *mockURLRepo) Get(id string) (domain.ShortenedURL, bool) {
	url, ok := r.data[id]
	return url, ok
}

func (r *mockURLRepo) Exists(id string) bool {
	_, ok := r.data[id]
	return ok
}

type mockVisitRepo struct {
	data map[string][]domain.Visit
}

func newMockVisitRepo() *mockVisitRepo {
	return &mockVisitRepo{data: make(map[string][]domain.Visit)}
}

func (r *mockVisitRepo) Record(visit domain.Visit) {
	r.data[visit.ShortURLID] = append(r.data[visit.ShortURLID], visit)
}

func (r *mockVisitRepo) GetByURL(id string) []domain.Visit {
	return r.data[id]
}

type mockCache struct {
	data map[string]string
}

func newMockCache() *mockCache {
	return &mockCache{data: make(map[string]string)}
}

func (c *mockCache) Get(key string) (string, bool) {
	val, ok := c.data[key]
	return val, ok
}

func (c *mockCache) Set(key, value string) {
	c.data[key] = value
}

func TestCreateShortURL(t *testing.T) {
	urls := newMockURLRepo()
	visits := newMockVisitRepo()
	cache := newMockCache()
	svc := NewURLService(urls, visits, cache)

	key, err := svc.CreateShortURL("https://example.com", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(key) != 6 {
		t.Errorf("expected key length 6, got %d", len(key))
	}

	url, found := svc.GetURL(key)
	if !found {
		t.Fatal("expected URL to be found")
	}
	if url.OriginalURL != "https://example.com" {
		t.Errorf("expected original URL to be https://example.com, got %s", url.OriginalURL)
	}
}

func TestCreateShortURL_CustomKey(t *testing.T) {
	urls := newMockURLRepo()
	visits := newMockVisitRepo()
	cache := newMockCache()
	svc := NewURLService(urls, visits, cache)

	key, err := svc.CreateShortURL("https://example.com", "mykey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key != "mykey" {
		t.Errorf("expected key 'mykey', got %s", key)
	}

	_, err = svc.CreateShortURL("https://other.com", "mykey")
	if err == nil {
		t.Error("expected error for duplicate custom key")
	}
}

func TestResolve(t *testing.T) {
	urls := newMockURLRepo()
	visits := newMockVisitRepo()
	cache := newMockCache()
	svc := NewURLService(urls, visits, cache)

	key, _ := svc.CreateShortURL("https://example.com", "")

	resolved, err := svc.Resolve(key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved != "https://example.com" {
		t.Errorf("expected https://example.com, got %s", resolved)
	}

	_, err = svc.Resolve("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent key")
	}
}

func TestResolve_Cache(t *testing.T) {
	urls := newMockURLRepo()
	visits := newMockVisitRepo()
	cache := newMockCache()
	svc := NewURLService(urls, visits, cache)

	key, _ := svc.CreateShortURL("https://example.com", "")

	// First resolve - populates cache
	svc.Resolve(key)

	// Delete from URL store - should still resolve from cache
	urls.data = make(map[string]domain.ShortenedURL)

	resolved, err := svc.Resolve(key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved != "https://example.com" {
		t.Errorf("expected https://example.com from cache, got %s", resolved)
	}
}

func TestRecordVisit(t *testing.T) {
	urls := newMockURLRepo()
	visits := newMockVisitRepo()
	cache := newMockCache()
	svc := NewURLService(urls, visits, cache)

	key, _ := svc.CreateShortURL("https://example.com", "")

	svc.RecordVisit(key, "Mozilla/5.0")

	visitsList := svc.GetVisits(key)
	if len(visitsList) != 1 {
		t.Fatalf("expected 1 visit, got %d", len(visitsList))
	}
	if visitsList[0].UserAgent != "Mozilla/5.0" {
		t.Errorf("expected Mozilla/5.0, got %s", visitsList[0].UserAgent)
	}
	if visitsList[0].ShortURLID != key {
		t.Errorf("expected short URL ID %s, got %s", key, visitsList[0].ShortURLID)
	}
	if visitsList[0].Timestamp.After(time.Now()) {
		t.Error("visit timestamp should not be in the future")
	}
}
