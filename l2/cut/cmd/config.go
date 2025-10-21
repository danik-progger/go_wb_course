package cut

import (
	"flag"
	"fmt"
	"os"
)

type CutConfig struct {
	Fields    string
	Delimiter rune
	Separated bool
}

func ParseConfig() *CutConfig {
	config := &CutConfig{}
	var delimiterStr string

	flag.StringVar(&config.Fields, "f", "", "fields to show, e.g. '1,3-5'")
	flag.StringVar(&delimiterStr, "d", "	", "delimiter character")
	flag.BoolVar(&config.Separated, "s", false, "only lines with delimiter")

	flag.Parse()

	if config.Fields == "" {
		fmt.Fprintln(os.Stderr, "error: -f flag is required")
		os.Exit(1)
	}

	runes := []rune(delimiterStr)
	if len(runes) > 1 {
		fmt.Fprintln(os.Stderr, "error: delimiter must be a single character")
		os.Exit(1)
	}
	config.Delimiter = runes[0]

	return config
}
