package main

import (
	"api/queue"
	"api/routes"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	server := &http.Server{
		Addr:         ":3000",
		Handler:      routes.SetUpRoutes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Go(queue.ProcessQ)
}
