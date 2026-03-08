package grep

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// Helper to redirect stdout and capture it
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestProcessLines(t *testing.T) {
	lines := []string{
		"first line",
		"second line with match",
		"third line",
		"fourth line is another match",
		"fifth line",
		"sixth line",
	}

	testCases := []struct {
		name   string
		conf   *GrepConfig
		input  []string
		output string
	}{
		{
			name:   "Simple match",
			conf:   &GrepConfig{pattern: "match"},
			input:  lines,
			output: "second line with match\nfourth line is another match\n",
		},
		{
			name:   "Invert match",
			conf:   &GrepConfig{pattern: "match", invert: true},
			input:  lines,
			output: "first line\nthird line\nfifth line\nsixth line\n",
		},
		{
			name:   "Fixed match (should not interpret as regex)",
			conf:   &GrepConfig{pattern: "line.with", fixed: true},
			input:  []string{"line with", "line.with"},
			output: "line.with\n",
		},
		{
			name:   "Regex match",
			conf:   &GrepConfig{pattern: "line.with"},
			input:  []string{"line with", "line.with"},
			output: "line with\nline.with\n",
		},
		{
			name:   "Case-insensitive fixed match",
			conf:   &GrepConfig{pattern: "MATCH", fixed: true, ignoreCase: true},
			input:  lines,
			output: "second line with match\nfourth line is another match\n",
		},
		{
			name:   "Case-insensitive regex match",
			conf:   &GrepConfig{pattern: "MATCH", ignoreCase: true},
			input:  lines,
			output: "second line with match\nfourth line is another match\n",
		},
		{
			name:   "Count",
			conf:   &GrepConfig{pattern: "match", count: true},
			input:  lines,
			output: "2\n",
		},
		{
			name:   "Line numbers",
			conf:   &GrepConfig{pattern: "match", lineNum: true},
			input:  lines,
			output: "2:second line with match\n4:fourth line is another match\n",
		},
		{
			name:   "Before context",
			conf:   &GrepConfig{pattern: "match", before: 1, lineNum: true},
			input:  lines,
			output: "1-first line\n2:second line with match\n--\n3-third line\n4:fourth line is another match\n",
		},
		{
			name:   "After context",
			conf:   &GrepConfig{pattern: "match", after: 1, lineNum: true},
			input:  lines,
			output: "2:second line with match\n3-third line\n--\n4:fourth line is another match\n5-fifth line\n",
		},
		{
			name:   "Context (A and B)",
			conf:   &GrepConfig{pattern: "match", after: 1, before: 1, lineNum: true},
			input:  lines,
			output: "1-first line\n2:second line with match\n3-third line\n4:fourth line is another match\n5-fifth line\n",
		},
		{
			name:   "Context (C)",
			conf:   &GrepConfig{pattern: "match", context: 1, lineNum: true},
			input:  lines,
			output: "1-first line\n2:second line with match\n3-third line\n4:fourth line is another match\n5-fifth line\n",
		},
		{
			name:   "First line match with After",
			conf:   &GrepConfig{pattern: "first", after: 1, lineNum: true},
			input:  lines,
			output: "1:first line\n2-second line with match\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.conf.context > 0 {
				tc.conf.after = tc.conf.context
				tc.conf.before = tc.conf.context
			}

			output := captureOutput(func() {
				processLines(tc.input, tc.conf)
			})

			if output != tc.output {
				t.Errorf("expected:\n%q\ngot:\n%q", tc.output, output)
			}
		})
	}
}

func TestGrepReader(t *testing.T) {
	conf := &GrepConfig{pattern: "hello"}
	input := "hello world\n"
	reader := strings.NewReader(input)

	output := captureOutput(func() {
		Grep(conf, reader)
	})

	expected := "hello world\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}
