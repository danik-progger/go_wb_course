package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wget/downloader"
	"wget/fetcher"
)

func main() {
	rootURL := flag.String("url", "", "The root URL to start crawling from")
	depth := flag.Int("depth", 2, "Recursion depth for crawling")
	outputDir := flag.String("output", "output", "Directory to save downloaded files")
	workers := flag.Int("workers", 10, "Maximum number of concurrent workers")
	timeout := flag.Duration("timeout", 10*time.Second, "HTTP request timeout")
	flag.Parse()

	if *rootURL == "" {
		log.Fatal("url flag is required")
	}

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("\nReceived interrupt signal, shutting down...")
		cancel()
	}()

	f := fetcher.NewRealFetcher(*timeout)
	d := downloader.NewDownloader(f, *outputDir, *workers)

	if err := d.Run(ctx, *rootURL, *depth); err != nil {
		log.Fatalf("Download failed: %v", err)
	}

	log.Println("Download complete.")
}
