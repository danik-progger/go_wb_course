package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultTimeout = 10 * time.Second

type Fetcher interface {
	Fetch(ctx context.Context, url string) (io.ReadCloser, error)
}

type realFetcher struct {
	client *http.Client
}

func NewRealFetcher(timeout time.Duration) *realFetcher {
	return &realFetcher{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (f *realFetcher) Fetch(ctx context.Context, urlString string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlString, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request for %s: %w", urlString, err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", urlString, err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("fetching %s: %s", urlString, resp.Status)
	}

	return resp.Body, nil
}
