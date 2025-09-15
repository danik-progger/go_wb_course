package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"l0/server"
)

func main() {
	port := ":8080"
	srv := server.InitServer(port)

	httpServer := &http.Server{
		Addr:    port,
		Handler: srv.Router,
	}

	go func() {
		fmt.Printf("🟢 Server starting on %s\n", srv.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("🔴 Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🟢 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("🔴 Server shutdown failed:", err)
	}

	log.Println("🟢 Server exited properly")
}
