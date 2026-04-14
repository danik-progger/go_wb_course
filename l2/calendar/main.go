package main

import (
	"calendar/infra/server"
	"flag"
)

func main() {
	port := flag.String("port", "", "server port (or use PORT env var)")
	flag.Parse()
	s := server.NewServer(*port)
	s.Run()
}
