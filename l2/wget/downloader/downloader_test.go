package downloader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

type mockFetcher struct {
	responses map[string]*mockResponse
	mu        sync.Mutex
	callCount map[string]int
}

type mockResponse struct {
	body string
	err  error
}

func newMockFetcher() *mockFetcher {
	return &mockFetcher{
		responses: make(map[string]*mockResponse),
		callCount: make(map[string]int),
	}
}

func (m *mockFetcher) addResponse(urlString, body string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responses[urlString] = &mockResponse{body: body}
}

func (m *mockFetcher) addError(urlString string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responses[urlString] = &mockResponse{err: err}
}

func (m *mockFetcher) Fetch(ctx context.Context, urlString string) (io.ReadCloser, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callCount[urlString]++

	resp, ok := m.responses[urlString]
	if !ok {
		return nil, fmt.Errorf("no mock response for URL: %s", urlString)
	}

	if resp.err != nil {
		return nil, resp.err
	}

	return io.NopCloser(strings.NewReader(resp.body)), nil
}

func (m *mockFetcher) getCallCount(urlString string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.callCount[urlString]
}

func createTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func readFileContent(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}
	return string(content)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func TestNewDownloader(t *testing.T) {
	fetcher := newMockFetcher()
	outputDir := "/tmp/test"
	maxWorkers := 5

	d := NewDownloader(fetcher, outputDir, maxWorkers)

	if d.fetcher != fetcher {
		t.Error("Fetcher not set correctly")
	}
	if d.outputDir != outputDir {
		t.Errorf("Expected outputDir %s, got %s", outputDir, d.outputDir)
	}
	if d.maxWorkers != maxWorkers {
		t.Errorf("Expected maxWorkers %d, got %d", maxWorkers, d.maxWorkers)
	}
	if d.visited == nil {
		t.Error("Visited map not initialized")
	}
	if d.sem == nil {
		t.Error("Semaphore not initialized")
	}
}

func TestDownloader_Run_SimplePage(t *testing.T) {
	fetcher := newMockFetcher()
	fetcher.addResponse("http://example.com/", `
		<html>
		<head><title>Test</title></head>
		<body>
			<h1>Hello World</h1>
			<a href="/about">About</a>
		</body>
		</html>
	`)

	outputDir := createTestDir(t)
	d := NewDownloader(fetcher, outputDir, 2)

	ctx := context.Background()
	err := d.Run(ctx, "http://example.com/", 1)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	expectedPath := filepath.Join(outputDir, "example.com", "index.html")
	if !fileExists(expectedPath) {
		t.Errorf("Expected file %s to exist", expectedPath)
	}

	content := readFileContent(t, expectedPath)
	if !strings.Contains(content, "Hello World") {
		t.Errorf("Expected content to contain 'Hello World', got:\n%s", content)
	}
}

func TestDownloader_Run_MultiplePages(t *testing.T) {
	fetcher := newMockFetcher()
	fetcher.addResponse("http://example.com/", `
		<html>
		<body>
			<a href="/page1">Page 1</a>
			<a href="/page2">Page 2</a>
		</body>
		</html>
	`)
	fetcher.addResponse("http://example.com/page1", `<html><body>Page 1</body></html>`)
	fetcher.addResponse("http://example.com/page2", `<html><body>Page 2</body></html>`)

	outputDir := createTestDir(t)
	d := NewDownloader(fetcher, outputDir, 2)

	ctx := context.Background()
	err := d.Run(ctx, "http://example.com/", 2)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	files := []string{
		filepath.Join(outputDir, "example.com", "index.html"),
		filepath.Join(outputDir, "example.com", "page1"),
		filepath.Join(outputDir, "example.com", "page2"),
	}

	for _, f := range files {
		if !fileExists(f) {
			t.Errorf("Expected file %s to exist", f)
		}
	}
}

func TestDownloader_Run_DepthLimit(t *testing.T) {
	fetcher := newMockFetcher()
	fetcher.addResponse("http://example.com/", `
		<html>
		<body>
			<a href="/level1">Level 1</a>
		</body>
		</html>
	`)
	fetcher.addResponse("http://example.com/level1", `
		<html>
		<body>
			<a href="/level2">Level 2</a>
		</body>
		</html>
	`)

	outputDir := createTestDir(t)
	d := NewDownloader(fetcher, outputDir, 2)

	ctx := context.Background()
	err := d.Run(ctx, "http://example.com/", 2)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	rootPath := filepath.Join(outputDir, "example.com", "index.html")
	level1Path := filepath.Join(outputDir, "example.com", "level1")
	level2Path := filepath.Join(outputDir, "example.com", "level2")

	if !fileExists(rootPath) {
		t.Errorf("Expected root file %s to exist", rootPath)
	}

	if !fileExists(level1Path) {
		t.Errorf("Expected level1 file %s to exist", level1Path)
	}

	if fileExists(level2Path) {
		t.Errorf("Expected level2 file %s to NOT exist (depth limit exceeded)", level2Path)
	}
}

func TestDownloader_Run_FetchError(t *testing.T) {
	fetcher := newMockFetcher()
	fetcher.addError("http://example.com/", fmt.Errorf("network error"))

	outputDir := createTestDir(t)
	d := NewDownloader(fetcher, outputDir, 2)

	ctx := context.Background()
	err := d.Run(ctx, "http://example.com/", 1)
	if err != nil {
		t.Fatalf("Run should not return error for fetch failure: %v", err)
	}

	expectedPath := filepath.Join(outputDir, "example.com", "index.html")
	if fileExists(expectedPath) {
		t.Errorf("Expected file %s to NOT exist after fetch error", expectedPath)
	}
}

func TestDownloader_Run_AvoidsDuplicateVisits(t *testing.T) {
	fetcher := newMockFetcher()
	fetcher.addResponse("http://example.com/", `
		<html>
		<body>
			<a href="/page1">Page 1</a>
			<a href="/page2">Page 2</a>
		</body>
		</html>
	`)
	fetcher.addResponse("http://example.com/page1", `
		<html>
		<body>
			<a href="/page2">Page 2</a>
		</body>
		</html>
	`)
	fetcher.addResponse("http://example.com/page2", `<html><body>Page 2</body></html>`)

	outputDir := createTestDir(t)
	d := NewDownloader(fetcher, outputDir, 2)

	ctx := context.Background()
	err := d.Run(ctx, "http://example.com/", 2)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	page2Count := fetcher.getCallCount("http://example.com/page2")
	if page2Count != 1 {
		t.Errorf("Expected page2 to be fetched exactly 1 time, got %d", page2Count)
	}
}

func TestDownloader_Run_SameDomainOnly(t *testing.T) {
	fetcher := newMockFetcher()
	fetcher.addResponse("http://example.com/", `
		<html>
		<body>
			<a href="http://example.com/local">Local</a>
			<a href="http://other.com/remote">Remote</a>
		</body>
		</html>
	`)
	fetcher.addResponse("http://example.com/local", `<html><body>Local</body></html>`)

	outputDir := createTestDir(t)
	d := NewDownloader(fetcher, outputDir, 2)

	ctx := context.Background()
	err := d.Run(ctx, "http://example.com/", 2)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	localPath := filepath.Join(outputDir, "example.com", "local")
	if !fileExists(localPath) {
		t.Errorf("Expected local file %s to exist", localPath)
	}

	if fetcher.getCallCount("http://other.com/remote") > 0 {
		t.Error("Should not fetch URLs from different domains")
	}
}

func TestDownloader_Run_LinkRewriting(t *testing.T) {
	fetcher := newMockFetcher()
	fetcher.addResponse("http://example.com/", `
		<html>
		<body>
			<a href="/about">About</a>
			<a href="http://example.com/contact">Contact</a>
		</body>
		</html>
	`)

	outputDir := createTestDir(t)
	d := NewDownloader(fetcher, outputDir, 2)

	ctx := context.Background()
	err := d.Run(ctx, "http://example.com/", 1)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	expectedPath := filepath.Join(outputDir, "example.com", "index.html")
	content := readFileContent(t, expectedPath)

	if !strings.Contains(content, `href="about"`) {
		t.Errorf("Expected link to be rewritten to relative path 'about', got:\n%s", content)
	}

	if !strings.Contains(content, `href="contact"`) {
		t.Errorf("Expected link to be rewritten to relative path 'contact', got:\n%s", content)
	}
}

func TestDownloader_getFilePath(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "root path becomes index.html",
			url:      "http://example.com/",
			expected: "example.com/index.html",
		},
		{
			name:     "empty path becomes index.html",
			url:      "http://example.com",
			expected: "example.com/index.html",
		},
		{
			name:     "specific path",
			url:      "http://example.com/about",
			expected: "example.com/about",
		},
		{
			name:     "nested path",
			url:      "http://example.com/docs/readme",
			expected: "example.com/docs/readme",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDownloader(newMockFetcher(), "/output", 1)
			u, err := url.Parse(tt.url)
			if err != nil {
				t.Fatalf("Failed to parse URL: %v", err)
			}

			result := d.getFilePath(u)
			expected := filepath.Join("/output", tt.expected)

			if result != expected {
				t.Errorf("Expected %s, got %s", expected, result)
			}
		})
	}
}

func TestDownloader_saveFile(t *testing.T) {
	outputDir := createTestDir(t)
	d := NewDownloader(newMockFetcher(), outputDir, 1)

	content := strings.NewReader("Test file content")
	testPath := filepath.Join(outputDir, "test", "nested", "file.txt")

	err := d.saveFile(testPath, content)
	if err != nil {
		t.Fatalf("saveFile failed: %v", err)
	}

	if !fileExists(testPath) {
		t.Errorf("Expected file %s to exist", testPath)
	}

	fileContent := readFileContent(t, testPath)
	if fileContent != "Test file content" {
		t.Errorf("Expected 'Test file content', got '%s'", fileContent)
	}
}

func TestDownloader_isSameDomain(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		testURL  string
		expected bool
	}{
		{
			name:     "same domain",
			baseURL:  "http://example.com/page",
			testURL:  "http://example.com/other",
			expected: true,
		},
		{
			name:     "different domain",
			baseURL:  "http://example.com/page",
			testURL:  "http://other.com/page",
			expected: false,
		},
		{
			name:     "same domain different scheme",
			baseURL:  "http://example.com/page",
			testURL:  "https://example.com/page",
			expected: true,
		},
		{
			name:     "invalid URL",
			baseURL:  "http://example.com/page",
			testURL:  "://invalid-url",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDownloader(newMockFetcher(), "/output", 1)
			baseURL, err := url.Parse(tt.baseURL)
			if err != nil {
				t.Fatalf("Failed to parse base URL: %v", err)
			}

			result := d.isSameDomain(tt.testURL, baseURL)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDownloader_markVisited(t *testing.T) {
	d := NewDownloader(newMockFetcher(), "/output", 1)

	url1 := "http://example.com/page1"
	url2 := "http://example.com/page2"

	if !d.markVisited(url1) {
		t.Error("First visit should return true")
	}

	if d.markVisited(url1) {
		t.Error("Second visit to same URL should return false")
	}

	if !d.markVisited(url2) {
		t.Error("Visit to different URL should return true")
	}
}

func TestDownloader_Run_ContextCancellation(t *testing.T) {
	fetcher := newMockFetcher()
	fetcher.addResponse("http://example.com/", `
		<html>
		<body>
			<a href="/page1">Page 1</a>
		</body>
		</html>
	`)

	outputDir := createTestDir(t)
	d := NewDownloader(fetcher, outputDir, 2)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := d.Run(ctx, "http://example.com/", 2)
	if err != nil {
		t.Fatalf("Run should handle cancelled context gracefully: %v", err)
	}
}

func TestDownloader_Run_ConcurrentAccess(t *testing.T) {
	fetcher := newMockFetcher()

	for i := 0; i < 10; i++ {
		fetcher.addResponse(
			fmt.Sprintf("http://example.com/page%d", i),
			fmt.Sprintf(`<html><body><a href="/page%d">Next</a></body></html>`, (i+1)%10),
		)
	}

	outputDir := createTestDir(t)
	d := NewDownloader(fetcher, outputDir, 5)

	ctx := context.Background()
	err := d.Run(ctx, "http://example.com/page0", 2)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if len(d.visited) == 0 {
		t.Error("Expected some pages to be visited")
	}
}

func TestDownloader_Run_EmptyResponse(t *testing.T) {
	fetcher := newMockFetcher()
	fetcher.addResponse("http://example.com/", "")

	outputDir := createTestDir(t)
	d := NewDownloader(fetcher, outputDir, 2)

	ctx := context.Background()
	err := d.Run(ctx, "http://example.com/", 1)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	expectedPath := filepath.Join(outputDir, "example.com", "index.html")
	if !fileExists(expectedPath) {
		t.Errorf("Expected file %s to exist even for empty response", expectedPath)
	}
}

func TestDownloader_Run_HTMLWithSpecialChars(t *testing.T) {
	fetcher := newMockFetcher()
	fetcher.addResponse("http://example.com/", `
		<html>
		<body>
			<h1>Special chars: &lt;&gt;&amp;"'</h1>
			<a href="/p%C3%A1ge">Página</a>
		</body>
		</html>
	`)

	outputDir := createTestDir(t)
	d := NewDownloader(fetcher, outputDir, 2)

	ctx := context.Background()
	err := d.Run(ctx, "http://example.com/", 1)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	expectedPath := filepath.Join(outputDir, "example.com", "index.html")
	content := readFileContent(t, expectedPath)

	if !strings.Contains(content, "Special chars:") {
		t.Errorf("Expected content to handle special characters, got:\n%s", content)
	}
}

func TestDownloader_Run_MultipleWorkers(t *testing.T) {
	fetcher := newMockFetcher()

	fetcher.addResponse("http://example.com/", `
		<html>
		<body>
			<a href="/page1">1</a>
			<a href="/page2">2</a>
			<a href="/page3">3</a>
			<a href="/page4">4</a>
		</body>
		</html>
	`)

	for i := 1; i <= 4; i++ {
		fetcher.addResponse(
			fmt.Sprintf("http://example.com/page%d", i),
			fmt.Sprintf(`<html><body>Page %d</body></html>`, i),
		)
	}

	outputDir := createTestDir(t)
	d := NewDownloader(fetcher, outputDir, 10)

	ctx := context.Background()
	err := d.Run(ctx, "http://example.com/", 2)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for i := 1; i <= 4; i++ {
		path := filepath.Join(outputDir, "example.com", fmt.Sprintf("page%d", i))
		if !fileExists(path) {
			t.Errorf("Expected %s to exist", path)
		}
	}
}

func TestDownloader_saveFile_DirectoryCreationError(t *testing.T) {
	outputDir := createTestDir(t)
	d := NewDownloader(newMockFetcher(), outputDir, 1)

	invalidPath := filepath.Join(outputDir, string([]byte{0}), "file.txt")
	content := bytes.NewReader([]byte("test"))

	err := d.saveFile(invalidPath, content)
	if err == nil {
		t.Error("Expected error for invalid path")
	}
}

func TestDownloader_Run_OnlyRootAtDepthZero(t *testing.T) {
	fetcher := newMockFetcher()
	fetcher.addResponse("http://example.com/", `<html><body>Root</body></html>`)

	outputDir := createTestDir(t)
	d := NewDownloader(fetcher, outputDir, 2)

	ctx := context.Background()
	err := d.Run(ctx, "http://example.com/", 0)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	expectedPath := filepath.Join(outputDir, "example.com", "index.html")
	if fileExists(expectedPath) {
		t.Errorf("Expected no files to be downloaded at depth 0")
	}
}
