package grep

import (
	"flag"
	"fmt"
	"os"
)

type GrepConfig struct {
	after      int
	before     int
	context    int
	count      bool
	ignoreCase bool
	invert     bool
	fixed      bool
	lineNum    bool
	pattern    string
	filePath   string
}

func NewGrepConfig() *GrepConfig {
	config := &GrepConfig{}

	flag.IntVar(&config.after, "A", 0, "print N lines of trailing context after matching lines")
	flag.IntVar(&config.before, "B", 0, "print N lines of leading context before matching lines")
	flag.IntVar(&config.context, "C", 0, "print N lines of leading and trailing context")
	flag.BoolVar(&config.count, "c", false, "print only a count of selected lines")
	flag.BoolVar(&config.ignoreCase, "i", false, "ignore case distinctions")
	flag.BoolVar(&config.invert, "v", false, "select non-matching lines")
	flag.BoolVar(&config.fixed, "F", false, "interpret pattern as a fixed string")
	flag.BoolVar(&config.lineNum, "n", false, "prefix each line of output with the line number")

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: grep [options] pattern [file]")
		os.Exit(1)
	}

	// If -C is used, it overrides -A and -B
	if config.context > 0 {
		config.after = config.context
		config.before = config.context
	}

	config.pattern = args[0]
	if len(args) > 1 {
		config.filePath = args[1]
	}

	return config
}
