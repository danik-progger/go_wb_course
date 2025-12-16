package types

// WorkItem represents a unit of work for grep
type WorkItem struct {
	Line       string
	LineNumber int
}

// SearchRequest represents a request to search for a pattern in a file
type SearchRequest struct {
	Pattern string `json:"pattern"`
	File    string `json:"file"`
	ID      string `json:"id"`
}

// SearchResponse represents the result of a search operation
type SearchResponse struct {
	ID      string   `json:"id"`
	Results []string `json:"results"`
	Error   string   `json:"error,omitempty"`
	Server  string   `json:"server"`
}

// SearchResult represents the combined result from multiple servers
type SearchResult struct {
	ID       string              `json:"id"`
	Results  map[string][]string `json:"results"` // server -> results
	Servers  []string            `json:"servers"`
	Quorum   int                 `json:"quorum"`
	Complete bool                `json:"complete"`
}
