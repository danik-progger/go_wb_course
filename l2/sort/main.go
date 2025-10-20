package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type sorter struct {
	lines    []string
	column   int
	numeric  bool
	reverse  bool
	unique   bool
	month    bool
	ignore   bool
	check    bool
	human    bool
}

var monthMap = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6,
	"jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12,
}

func main() {
	k := flag.Int("k", 1, "sort via a key; the key is a column number")
	n := flag.Bool("n", false, "compare according to string numerical value")
	r := flag.Bool("r", false, "reverse the result of comparisons")
	u := flag.Bool("u", false, "output only the first of an equal run")
	M := flag.Bool("M", false, "compare (unknown) < 'JAN' < ... < 'DEC'")
	b := flag.Bool("b", false, "ignore leading blanks")
	c := flag.Bool("c", false, "check for sorted input; do not sort")
	h := flag.Bool("h", false, "compare human readable numbers (e.g., 2K 1G)")

	flag.Parse()

	var reader io.Reader
	if flag.NArg() > 0 {
		file, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			os.Exit(1)
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	lines, err := readLines(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	s := &sorter{
		lines:    lines,
		column:   *k,
		numeric:  *n,
		reverse:  *r,
		unique:   *u,
		month:    *M,
		ignore:   *b,
		check:    *c,
		human:    *h,
	}

	if s.check {
		if !s.isSorted() {
			fmt.Println("sort: input is not sorted")
			os.Exit(1)
		}
		return
	}

	s.sort()

	if s.unique {
		s.lines = s.removeDuplicates()
	}

	for _, line := range s.lines {
		fmt.Println(line)
	}
}

func readLines(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func (s *sorter) sort() {
	sort.SliceStable(s.lines, func(i, j int) bool {
		a, b := s.lines[i], s.lines[j]
		res := s.compare(a, b)
		if s.reverse {
			return !res
		}
		return res
	})
}

func (s *sorter) isSorted() bool {
	for i := 1; i < len(s.lines); i++ {
		if !s.compare(s.lines[i-1], s.lines[i]) {
			return false
		}
	}
	return true
}

func (s *sorter) compare(a, b string) bool {
	valA := s.getVal(a)
	valB := s.getVal(b)

	if s.numeric {
		numA, errA := strconv.ParseFloat(valA, 64)
		numB, errB := strconv.ParseFloat(valB, 64)
		if errA == nil && errB == nil {
			return numA < numB
		}
	}

	if s.month {
		monthA, okA := monthMap[strings.ToLower(valA)]
		monthB, okB := monthMap[strings.ToLower(valB)]
		if okA && okB {
			return monthA < monthB
		}
	}

	if s.human {
		numA, errA := parseHuman(valA)
		numB, errB := parseHuman(valB)
		if errA == nil && errB == nil {
			return numA < numB
		}
	}

	return valA < valB
}

func (s *sorter) getVal(line string) string {
	if s.ignore {
		line = strings.TrimSpace(line)
	}

	fields := strings.Fields(line)
	if s.column > 0 && s.column <= len(fields) {
		return fields[s.column-1]
	}
	return line
}

func (s *sorter) removeDuplicates() []string {
	if len(s.lines) == 0 {
		return s.lines
	}
	result := []string{s.lines[0]}
	for i := 1; i < len(s.lines); i++ {
		if s.lines[i] != s.lines[i-1] {
			result = append(result, s.lines[i])
		}
	}
	return result
}

func parseHuman(s string) (int64, error) {
	s = strings.ToUpper(s)
	var multiplier int64 = 1
	lastChar := s[len(s)-1]

	switch lastChar {
	case 'K':
		multiplier = 1024
		s = s[:len(s)-1]
	case 'M':
		multiplier = 1024 * 1024
		s = s[:len(s)-1]
	case 'G':
		multiplier = 1024 * 1024 * 1024
		s = s[:len(s)-1]
	}

	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return val * multiplier, nil
}
