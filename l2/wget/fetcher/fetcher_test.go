package fetcher

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type mockFetcher struct {
	responses map[string]string
}

func (f *mockFetcher) Fetch(ctx context.Context, urlString string) (io.ReadCloser, error) {
	body, ok := f.responses[urlString]
	if !ok {
		return nil, &httpError{urlString, "not found"}
	}
	return io.NopCloser(strings.NewReader(body)), nil
}

type httpError struct {
	url string
	msg string
}

func (e *httpError) Error() string {
	return e.msg + ": " + e.url
}

func newMockFetcher(responses map[string]string) *mockFetcher {
	return &mockFetcher{responses: responses}
}

func TestRealFetcherTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	fetcher := NewRealFetcher(50 * time.Millisecond)

	ctx := context.Background()
	_, err := fetcher.Fetch(ctx, server.URL)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestRealFetcherSuccess(t *testing.T) {
	expectedContent := "Hello, World!"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expectedContent))
	}))
	defer server.Close()

	fetcher := NewRealFetcher(5 * time.Second)

	ctx := context.Background()
	body, err := fetcher.Fetch(ctx, server.URL)
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	defer body.Close()

	content, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	if string(content) != expectedContent {
		t.Errorf("Content = %q, want %q", string(content), expectedContent)
	}
}

func TestRealFetcherNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	fetcher := NewRealFetcher(5 * time.Second)

	ctx := context.Background()
	_, err := fetcher.Fetch(ctx, server.URL)
	if err == nil {
		t.Error("Expected error for 404, got nil")
	}
}
