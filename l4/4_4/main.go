package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof" // Import pprof for profiling
	"runtime/debug"

	"profiler/handlers"
)

var (
	port      = flag.String("port", "8080", "Port to run the server on")
	gcPercent = flag.Int("gcpercent", 100, "GOGC value (percentage of heap growth that triggers GC)")
)

func main() {
	flag.Parse()

	// Set GC percent if specified
	if *gcPercent != 100 {
		debug.SetGCPercent(*gcPercent)
		log.Printf("GC percent set to: %d", *gcPercent)
	}

	// Initialize metrics
	handlers.InitMetrics()

	// Start metrics collection in a goroutine
	go handlers.CollectMetrics()

	// Register HTTP handlers
	http.HandleFunc("/metrics", handlers.MetricsHandler)
	http.HandleFunc("/health", handlers.HealthHandler)

	// pprof endpoints are automatically registered by importing net/http/pprof

	// Create a sample endpoint that allocates memory to demonstrate the metrics
	http.HandleFunc("/allocate", handlers.AllocateMemoryHandler)

	addr := ":" + *port
	log.Printf("Server starting on http://localhost%s", addr)
	log.Printf("Metrics available at http://localhost%s/metrics", addr)
	log.Printf("Health check at http://localhost%s/health", addr)
	log.Printf("PProf available at http://localhost%s/debug/pprof/", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}
