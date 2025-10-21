package main

import (
	grep "grep/cmd"
	"io"
)

func main() {
	config := grep.NewGrepConfig()
	var reader io.Reader
	grep.Grep(config, reader)
}
