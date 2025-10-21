package main

import (
	cut "cut/cmd"
	"io"
)

func main() {
	config := cut.ParseConfig()
	var reader io.Reader
	cut.Cut(config, reader)
}
