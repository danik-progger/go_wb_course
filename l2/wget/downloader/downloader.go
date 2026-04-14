package downloader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"wget/fetcher"
	"wget/parser"
)

type Downloader struct {
	fetcher    fetcher.Fetcher
	outputDir  string
	maxWorkers int

	visited map[string]bool
	mu      sync.Mutex
	wg      sync.WaitGroup
	sem     chan struct{}
}

func NewDownloader(fetcher fetcher.Fetcher, outputDir string, maxWorkers int) *Downloader {
	return &Downloader{
		fetcher:    fetcher,
		outputDir:  outputDir,
		maxWorkers: maxWorkers,
		visited:    make(map[string]bool),
		sem:        make(chan struct{}, maxWorkers),
	}
}

func (d *Downloader) Run(ctx context.Context, rootURL string, depth int) error {
	d.wg.Add(1)
	go d.crawl(ctx, rootURL, depth)
	d.wg.Wait()
	return nil
}

func (d *Downloader) crawl(ctx context.Context, urlString string, depth int) {
	defer d.wg.Done()

	if depth <= 0 {
		return
	}

	if !d.markVisited(urlString) {
		return
	}

	fmt.Printf("Crawling: %s at depth %d\n", urlString, depth)

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		log.Printf("Failed to parse URL %s: %v", urlString, err)
		return
	}

	body, err := d.fetcher.Fetch(ctx, urlString)
	if err != nil {
		log.Printf("Failed to fetch %s: %v", urlString, err)
		return
	}
	defer body.Close()

	var buf bytes.Buffer
	tee := io.TeeReader(body, &buf)

	rewriter := parser.NewLinkRewriter(parsedURL, d.outputDir)
	links, rewrittenContent, err := rewriter.ExtractAndRewriteLinks(tee)
	if err != nil {
		log.Printf("Failed to extract links from %s: %v", urlString, err)
		return
	}

	filePath := d.getFilePath(parsedURL)
	if err := d.saveFile(filePath, rewrittenContent); err != nil {
		log.Printf("Failed to save file %s: %v", filePath, err)
		return
	}

	for _, link := range links {
		if !d.isSameDomain(link, parsedURL) {
			continue
		}

		d.wg.Add(1)
		d.sem <- struct{}{}
		go func(link string) {
			defer func() { <-d.sem }()
			d.crawl(ctx, link, depth-1)
		}(link)
	}
}

func (d *Downloader) markVisited(urlString string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.visited[urlString] {
		return false
	}
	d.visited[urlString] = true
	return true
}

func (d *Downloader) getFilePath(u *url.URL) string {
	path := filepath.Join(d.outputDir, u.Host, u.Path)
	if strings.HasSuffix(u.Path, "/") || u.Path == "" {
		path = filepath.Join(path, "index.html")
	}
	return path
}

func (d *Downloader) saveFile(path string, content io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating directory for %s: %w", path, err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", path, err)
	}
	defer f.Close()

	if _, err := io.Copy(f, content); err != nil {
		return fmt.Errorf("writing content to %s: %w", path, err)
	}

	return nil
}

func (d *Downloader) isSameDomain(urlString string, baseURL *url.URL) bool {
	u, err := url.Parse(urlString)
	if err != nil {
		return false
	}
	return u.Host == baseURL.Host
}
