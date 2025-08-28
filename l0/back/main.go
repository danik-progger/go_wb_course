package main

import (
	"fmt"
	"log"
	"net/http"

	"l0/server"
)

func main() {
	port := ":8080"
	srv := server.InitServer(port)

	fmt.Printf("🟢 Server starting on %s\n", srv.Port)
	log.Fatal(http.ListenAndServe(srv.Port, srv.Router))
}
