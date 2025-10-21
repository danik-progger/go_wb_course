package cut

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func Cut(config *CutConfig, reader io.Reader) {
	fields, err := NewFields(config.Fields)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing fields:", err)
		os.Exit(1)
	}

	delimiterStr := string(config.Delimiter)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.Contains(line, delimiterStr) {
			if !config.Separated {
				fmt.Println(line)
			}
			continue
		}

		parts := strings.Split(line, delimiterStr)
		resultParts := fields.Extract(parts)

		fmt.Println(strings.Join(resultParts, delimiterStr))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}
}
