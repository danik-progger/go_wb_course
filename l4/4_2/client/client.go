package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"cli/types"
)

// Client represents a distributed grep client
type Client struct {
	servers []string
	quorum  int
	client  *http.Client
}

// NewClient creates a new Client instance
func NewClient(servers []string, quorum int) *Client {
	if quorum == 0 {
		quorum = len(servers)/2 + 1
	}

	return &Client{
		servers: servers,
		quorum:  quorum,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Search distributes the search request to multiple servers and returns results when quorum is reached
func (c *Client) Search(pattern, file string) ([]string, error) {
	// Create search request
	req := types.SearchRequest{
		Pattern: pattern,
		File:    file,
		ID:      fmt.Sprintf("search-%d", time.Now().Unix()),
	}

	// Channel to receive responses
	responseCh := make(chan types.SearchResponse, len(c.servers))

	// Start search on all servers
	var wg sync.WaitGroup
	for _, server := range c.servers {
		wg.Add(1)
		go func(srv string) {
			defer wg.Done()
			resp, err := c.sendSearchRequest(srv, req)
			if err != nil {
				// Send error response
				responseCh <- types.SearchResponse{
					ID:     req.ID,
					Error:  err.Error(),
					Server: srv,
				}
				return
			}
			responseCh <- resp
		}(server)
	}

	// Start a goroutine to close responseCh when all requests are done
	go func() {
		wg.Wait()
		close(responseCh)
	}()

	// Collect responses and check for quorum
	results := make(map[string][]string) // server -> results
	successfulServers := 0
	var errors []string

	// Use a ticker to periodically check for quorum
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Create a timeout
	timeout := time.NewTimer(30 * time.Second)
	defer timeout.Stop()

	for {
		select {
		case resp, ok := <-responseCh:
			if !ok {
				// All responses received
				if successfulServers >= c.quorum {
					return c.combineResults(results), nil
				}
				// Not enough servers responded successfully
				return nil, fmt.Errorf("insufficient successful responses for quorum: %d/%d, errors: %v", successfulServers, len(c.servers), errors)
			}

			if resp.Error != "" {
				errors = append(errors, fmt.Sprintf("%s: %s", resp.Server, resp.Error))
			} else {
				results[resp.Server] = resp.Results
				successfulServers++
			}

			// Check if we have reached quorum
			if successfulServers >= c.quorum {
				return c.combineResults(results), nil
			}

		case <-ticker.C:
			// Periodically check if we have reached quorum
			if successfulServers >= c.quorum {
				return c.combineResults(results), nil
			}

		case <-timeout.C:
			// Timeout reached
			return nil, fmt.Errorf("timeout waiting for quorum (%d/%d servers responded successfully), errors: %v", successfulServers, len(c.servers), errors)
		}
	}
}

// sendSearchRequest sends a search request to a specific server
func (c *Client) sendSearchRequest(server string, req types.SearchRequest) (types.SearchResponse, error) {
	url := fmt.Sprintf("http://%s/search", server)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return types.SearchResponse{}, err
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return types.SearchResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return types.SearchResponse{}, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	var searchResp types.SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return types.SearchResponse{}, err
	}

	return searchResp, nil
}

// combineResults merges results from different servers
func (c *Client) combineResults(results map[string][]string) []string {
	finalResults := make(map[string]bool) // Use map to avoid duplicates

	for _, serverResults := range results {
		for _, result := range serverResults {
			finalResults[result] = true
		}
	}

	// Convert map back to slice
	var combined []string
	for result := range finalResults {
		combined = append(combined, result)
	}

	return combined
}
