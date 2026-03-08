package main

import (
	grep "grep/cmd"
)

func main() {
	config := grep.NewGrepConfig()
	grep.Grep(config, nil)
}
