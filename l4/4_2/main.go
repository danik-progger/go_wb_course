package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"cli/client"
	"cli/server"
	"cli/types"
)

var (
	pattern = flag.String("pattern", "", "Pattern to search for")
	file    = flag.String("file", "", "File to search in")
	workers = flag.Int("workers", 4, "Number of worker goroutines")
	port    = flag.String("port", "8080", "Port to run the server on")
	remote  = flag.Bool("remote", false, "Run in client mode to connect to remote servers")
	servers = flag.String("servers", "", "Comma-separated list of server addresses (host:port)")
	quorum  = flag.Int("quorum", 0, "Minimum number of servers required for quorum (default: N/2 + 1)")
)

func main() {
	flag.Parse()

	if *remote {
		// Run as client connecting to remote servers
		if *servers == "" {
			log.Fatal("Servers list is required in remote mode")
		}
		if *pattern == "" {
			log.Fatal("Pattern is required")
		}

		serverList := strings.Split(*servers, ",")
		if *quorum == 0 {
			*quorum = len(serverList)/2 + 1
		}

		client := client.NewClient(serverList, *quorum)
		results, err := client.Search(*pattern, *file)
		if err != nil {
			log.Fatal(err)
		}

		// Print results
		for _, result := range results {
			fmt.Println(result)
		}
	} else {
		// Run as server
		if *pattern == "" && *file == "" {
			// If no pattern and file specified, run as server only
			log.Printf("Starting server on port %s", *port)
			s := server.NewServer(*port)
			s.Start()
		} else {
			// If both pattern and file are specified, run local search
			if *pattern == "" || *file == "" {
				log.Fatal("Both pattern and file are required for local search")
			}

			// Local grep operation
			results := localGrep(*pattern, *file, *workers)
			for _, result := range results {
				fmt.Println(result)
			}
		}
	}
}

// localGrep performs local grep operation with concurrency
func localGrep(pattern, filename string, numWorkers int) []string {
	// Read the file and split into chunks for parallel processing
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(fileContent), "\n")

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

	return finalResults
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
