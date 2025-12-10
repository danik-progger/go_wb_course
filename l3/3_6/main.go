package main

import (
	"fmt"
	"log"
	"net/http"
	"salesTracker/api"
	"salesTracker/db"

	"github.com/go-chi/chi/v5"
)

func main() {
	// Initialize database
	salesDB, err := db.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize API server
	apiServer := api.NewServer(salesDB)

	// Initialize frontend server
	frontendServer := api.NewFrontendServer(salesDB)

	// Create a main router
	router := chi.NewRouter()

	// Add CORS middleware
	router.Use(corsMiddleware)

	// Mount API routes under /api to avoid conflicts with frontend routes
	router.Mount("/api", apiServer.Router)

	// Mount frontend routes for the root path
	// This ensures that when users go to the root, they get the frontend
	router.Mount("/", frontendServer.Router)

	// Start HTTP server
	port := ":8080"

	// Log the URL where to access the frontend
	fmt.Printf("Server starting on http://localhost%s\n", port)
	fmt.Printf("Access the frontend at: http://localhost%s\n", port)

	log.Fatal(http.ListenAndServe(port, router))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
