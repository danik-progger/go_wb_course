package grep

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func Grep(conf *GrepConfig, reader io.Reader) {
	if conf.filePath != "" {
		file, err := os.Open(conf.filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "🔴 Error opening file: %v", err)
			os.Exit(1)
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	scanner := bufio.NewScanner(reader)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v", err)
		os.Exit(1)
	}

	processLines(lines, conf)
}

func processLines(lines []string, conf *GrepConfig) {
	var matcher func(string) bool
	if conf.fixed {
		pattern := conf.pattern
		if conf.ignoreCase {
			pattern = strings.ToLower(pattern)
		}
		matcher = func(line string) bool {
			if conf.ignoreCase {
				line = strings.ToLower(line)
			}
			return strings.Contains(line, pattern)
		}
	} else {
		pattern := conf.pattern
		if conf.ignoreCase {
			pattern = "(?i)" + pattern
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid regular expression: %v", err)
			os.Exit(1)
		}
		matcher = func(line string) bool {
			return re.MatchString(line)
		}
	}

	matchCount := 0
	linesToPrint := make(map[int]bool)
	matchIndices := []int{}

	for i, line := range lines {
		match := matcher(line)
		if conf.invert {
			match = !match
		}
		if match {
			matchCount++
			matchIndices = append(matchIndices, i)
		}
	}

	if conf.count {
		fmt.Println(matchCount)
		return
	}

	for _, idx := range matchIndices {
		start := max(idx-conf.before, 0)
		end := min(idx+conf.after, len(lines)-1)
		for i := start; i <= end; i++ {
			linesToPrint[i] = true
		}
	}

	print(lines, linesToPrint, conf)

}

func print(lines []string, linesToPrint map[int]bool, conf *GrepConfig) {
	var printed bool
	for i := range len(lines) {
		if i > 0 {
			if printed && !linesToPrint[i-1] {
				fmt.Println("--")
			}
			if conf.lineNum {
				fmt.Print(strconv.Itoa(i+1) + ":")
			}
			fmt.Println(lines[i])
			printed = true
		}
	}
}
