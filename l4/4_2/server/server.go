package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"cli/types"
)

// Server represents a grep server
type Server struct {
	port   string
	router *http.ServeMux
}

// NewServer creates a new Server instance
func NewServer(port string) *Server {
	s := &Server{
		port:   port,
		router: http.NewServeMux(),
	}

	s.setupRoutes()
	return s
}

// setupRoutes configures HTTP routes
func (s *Server) setupRoutes() {
	s.router.HandleFunc("/search", s.handleSearch)
	s.router.HandleFunc("/health", s.handleHealth)
}

// handleSearch receives search requests and returns grep results
func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req types.SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Perform grep operation
	results, err := grep(req.Pattern, req.File)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create response
	resp := types.SearchResponse{
		ID:      req.ID,
		Results: results,
		Server:  fmt.Sprintf("localhost:%s", s.port),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleHealth returns server health status
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"server": fmt.Sprintf("localhost:%s", s.port),
	})
}

// grep performs grep operation on a file
func grep(pattern, filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var results []string
	scanner := bufio.NewScanner(file)
	lineNum := 1

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, pattern) {
			results = append(results, fmt.Sprintf("%d:%s", lineNum, line))
		}
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// Start starts the server
func (s *Server) Start() {
	addr := fmt.Sprintf(":%s", s.port)
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, s.router))
}

// ConcurrentGrep performs concurrent grep on chunks of a file
func ConcurrentGrep(pattern, filename string, numWorkers int) ([]string, error) {
	// Read the entire file
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Split content into lines
	lines := strings.Split(string(content), "\n")

	// Channel for work items
	work := make(chan types.WorkItem, len(lines))

	// Channel for results
	results := make(chan string, len(lines))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(pattern, work, results, &wg)
	}

	// Send work to workers
	go func() {
		defer close(work)
		for i, line := range lines {
			work <- types.WorkItem{
				Line:       line,
				LineNumber: i + 1,
			}
		}
	}()

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var finalResults []string
	for result := range results {
		finalResults = append(finalResults, result)
	}

	return finalResults, nil
}

// worker processes work items and sends results
func worker(pattern string, work <-chan types.WorkItem, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for item := range work {
		if strings.Contains(item.Line, pattern) {
			results <- fmt.Sprintf("%d:%s", item.LineNumber, item.Line)
		}
	}
}
