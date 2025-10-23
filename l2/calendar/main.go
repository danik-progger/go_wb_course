package main

import (
	"calendar/infra/server"
)

func main() {
	s := server.NewServer(":8080")
	s.Run()
}
