package grep

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
)

func Grep(conf *GrepConfig, reader io.Reader) {
	if reader == nil {
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

	if len(matchIndices) == 0 {
		return
	}

	isMatch := make(map[int]bool)
	for _, idx := range matchIndices {
		isMatch[idx] = true
	}

	if conf.before == 0 && conf.after == 0 {
		for _, idx := range matchIndices {
			if conf.lineNum {
				fmt.Printf("%d:%s\n", idx+1, lines[idx])
			} else {
				fmt.Println(lines[idx])
			}
		}
		return
	}

	// With context
	var blocks [][2]int
	for _, idx := range matchIndices {
		start := max(idx-conf.before, 0)
		end := min(idx+conf.after, len(lines)-1)
		blocks = append(blocks, [2]int{start, end})
	}

	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i][0] < blocks[j][0]
	})

	merged := [][2]int{blocks[0]}
	for i := 1; i < len(blocks); i++ {
		last := &merged[len(merged)-1]
		current := blocks[i]
		if current[0] <= last[1] { // Merge overlapping blocks
			if current[1] > last[1] {
				last[1] = current[1]
			}
		} else {
			merged = append(merged, current)
		}
	}

	for i, block := range merged {
		if i > 0 {
			fmt.Println("--")
		}
		for lineIdx := block[0]; lineIdx <= block[1]; lineIdx++ {
			if conf.lineNum {
				separator := "-"
				if isMatch[lineIdx] {
					separator = ":"
				}
				fmt.Printf("%d%s%s\n", lineIdx+1, separator, lines[lineIdx])
			} else {
				fmt.Println(lines[lineIdx])
			}
		}
	}
}
