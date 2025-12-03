package main

import (
	"fmt"
	"imageProcessing/api"
	"imageProcessing/service"
	"log"
	"net/http"
)

func main() {
	// Create image service with storage path
	imageService := service.NewImageService("./storage")

	// Create server with the image service
	server := api.NewServer(imageService)

	// Start the server
	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, server.Router))
}
