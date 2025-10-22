package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

var (
	visited = make(map[string]bool)
	wg      sync.WaitGroup
	mu      sync.Mutex
)

type Fetcher interface {
	Fetch(url string) (io.ReadCloser, error)
}

type realFetcher struct{}

func (f *realFetcher) Fetch(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("fetching %s: %s", url, resp.Status)
	}
	return resp.Body, nil
}

func main() {
	rootURL := flag.String("url", "", "The root URL to start crawling from")
	depth := flag.Int("depth", 2, "Recursion depth for crawling")
	outputDir := flag.String("output", "output", "Directory to save downloaded files")
	flag.Parse()

	if *rootURL == "" {
		log.Fatal("url flag is required")
	}

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	fetcher := &realFetcher{}
	wg.Add(1)
	go crawl(*rootURL, *depth, *outputDir, fetcher)

	wg.Wait()
	fmt.Println("Download complete.")
}

func crawl(urlString string, depth int, outputDir string, fetcher Fetcher) {
	defer wg.Done()

	if depth <= 0 {
		return
	}

	mu.Lock()
	if visited[urlString] {
		mu.Unlock()
		return
	}
	visited[urlString] = true
	mu.Unlock()

	fmt.Printf("Crawling: %s at depth %d", urlString, depth)

	body, err := fetcher.Fetch(urlString)
	if err != nil {
		log.Printf("Failed to fetch %s: %v", urlString, err)
		return
	}
	defer body.Close()

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		log.Printf("Failed to parse URL %s: %v", urlString, err)
		return
	}

	// Read content to a buffer so we can both save and parse it
	var buf bytes.Buffer
	tee := io.TeeReader(body, &buf)

	links, err := extractLinks(tee, parsedURL)
	if err != nil {
		log.Printf("Failed to extract links from %s: %v", urlString, err)
	}

	// Save the original content
	filePath := getFilePath(parsedURL, outputDir)
	if err := saveFile(filePath, &buf); err != nil {
		log.Printf("Failed to save file %s: %v", filePath, err)
		return
	}

	for _, link := range links {
		if isSameDomain(link, parsedURL) {
			wg.Add(1)
			go crawl(link, depth-1, outputDir, fetcher)
		}
	}
}

func extractLinks(body io.Reader, baseURL *url.URL) ([]string, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "a", "link":
				for i, a := range n.Attr {
					if a.Key == "href" {
						resolvedURL, err := resolveURL(baseURL, a.Val)
						if err != nil {
							log.Printf("Could not resolve URL %s: %v", a.Val, err)
							continue
						}
						links = append(links, resolvedURL)
						// We don't rewrite the link in this simplified version
						// but a full implementation would change a.Val to the local path
						n.Attr[i].Val = a.Val
					}
				}
			case "img", "script":
				for i, a := range n.Attr {
					if a.Key == "src" {
						resolvedURL, err := resolveURL(baseURL, a.Val)
						if err != nil {
							log.Printf("Could not resolve URL %s: %v", a.Val, err)
							continue
						}
						links = append(links, resolvedURL)
						n.Attr[i].Val = a.Val
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return links, nil
}

func resolveURL(base *url.URL, href string) (string, error) {
	rel, err := url.Parse(href)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(rel).String(), nil
}

func getFilePath(u *url.URL, outputDir string) string {
	path := filepath.Join(outputDir, u.Host, u.Path)
	if strings.HasSuffix(u.Path, "/") || u.Path == "" {
		path = filepath.Join(path, "index.html")
	}
	return path
}

func saveFile(path string, content io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, content)
	return err
}

func isSameDomain(urlString string, baseURL *url.URL) bool {
	u, err := url.Parse(urlString)
	if err != nil {
		return false
	}
	return u.Host == baseURL.Host
}
