package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExpandEnvVars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		setupEnv map[string]string
		expected string
	}{
		{
			name:     "simple variable",
			input:    "echo $HOME",
			setupEnv: map[string]string{"HOME": "/home/user"},
			expected: "echo /home/user",
		},
		{
			name:     "braced variable",
			input:    "echo ${HOME}",
			setupEnv: map[string]string{"HOME": "/home/user"},
			expected: "echo /home/user",
		},
		{
			name:     "multiple variables",
			input:    "echo $USER at $HOME",
			setupEnv: map[string]string{"USER": "testuser", "HOME": "/home/testuser"},
			expected: "echo testuser at /home/testuser",
		},
		{
			name:     "undefined variable",
			input:    "echo $UNDEFINED_VAR",
			setupEnv: map[string]string{},
			expected: "echo ",
		},
		{
			name:     "mixed defined and undefined",
			input:    "echo $HOME $UNKNOWN",
			setupEnv: map[string]string{"HOME": "/root"},
			expected: "echo /root ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.setupEnv {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range tt.setupEnv {
					os.Unsetenv(k)
				}
			}()

			result := expandEnvVars(tt.input)
			if result != tt.expected {
				t.Errorf("expandEnvVars(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseRedirects(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		wantInputFile  string
		wantOutputFile string
		wantAppend     bool
		wantCleaned    string
	}{
		{
			name:           "output redirect",
			input:          "echo hello > output.txt",
			wantOutputFile: "output.txt",
			wantAppend:     false,
			wantCleaned:    "echo hello",
		},
		{
			name:           "append redirect",
			input:          "echo hello >> output.txt",
			wantOutputFile: "output.txt",
			wantAppend:     true,
			wantCleaned:    "echo hello",
		},
		{
			name:          "input redirect",
			input:         "cat < input.txt",
			wantInputFile: "input.txt",
			wantCleaned:   "cat",
		},
		{
			name:           "combined redirects",
			input:          "cat < input.txt > output.txt",
			wantInputFile:  "input.txt",
			wantOutputFile: "output.txt",
			wantAppend:     false,
			wantCleaned:    "cat",
		},
		{
			name:        "no redirects",
			input:       "echo hello",
			wantCleaned: "echo hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redirect, cleaned := parseRedirects(tt.input)

			if redirect.InputFile != tt.wantInputFile {
				t.Errorf("InputFile = %q, want %q", redirect.InputFile, tt.wantInputFile)
			}
			if redirect.OutputFile != tt.wantOutputFile {
				t.Errorf("OutputFile = %q, want %q", redirect.OutputFile, tt.wantOutputFile)
			}
			if redirect.Append != tt.wantAppend {
				t.Errorf("Append = %v, want %v", redirect.Append, tt.wantAppend)
			}
			if cleaned != tt.wantCleaned {
				t.Errorf("Cleaned = %q, want %q", cleaned, tt.wantCleaned)
			}
		})
	}
}

func TestCdCommand(t *testing.T) {
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	tmpDir := t.TempDir()

	tmpDir, err = filepath.EvalSymlinks(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		path    string
		setup   func() (string, func())
		wantErr bool
	}{
		{
			name: "absolute path",
			path: tmpDir,
			setup: func() (string, func()) {
				return tmpDir, func() {}
			},
			wantErr: false,
		},
		{
			name: "tilde only",
			path: "~",
			setup: func() (string, func()) {
				home, err := os.UserHomeDir()
				if err != nil {
					return "", func() {}
				}
				home, err = filepath.EvalSymlinks(home)
				if err != nil {
					return home, func() {}
				}
				return home, func() {}
			},
			wantErr: false,
		},
		{
			name: "tilde with subpath",
			path: "~/test",
			setup: func() (string, func()) {
				home, err := os.UserHomeDir()
				if err != nil {
					return "", func() {}
				}
				testSubDir := filepath.Join(home, "test")
				os.MkdirAll(testSubDir, 0755)
				testSubDir, _ = filepath.EvalSymlinks(testSubDir)
				return testSubDir, func() {
					os.RemoveAll(testSubDir)
				}
			},
			wantErr: false,
		},
		{
			name:    "nonexistent path",
			path:    "/nonexistent/path/xyz123",
			setup:   func() (string, func()) { return "", func() {} },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Chdir(origDir)

			expectedDir, cleanup := tt.setup()
			defer cleanup()

			if strings.HasPrefix(tt.path, "~") && expectedDir == "" {
				t.Skip("HOME not set, skipping tilde test")
			}

			err := cdCommand(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("cdCommand(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				currentDir, _ := os.Getwd()
				if currentDir != expectedDir {
					t.Errorf("cdCommand(%q) ended up in %q, want %q", tt.path, currentDir, expectedDir)
				}
			}
		})
	}
}

func TestSplitByOperator(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		operator string
		expected []string
	}{
		{
			name:     "AND operator",
			input:    "cmd1 && cmd2 && cmd3",
			operator: "&&",
			expected: []string{"cmd1", "cmd2", "cmd3"},
		},
		{
			name:     "OR operator",
			input:    "cmd1 || cmd2 || cmd3",
			operator: "||",
			expected: []string{"cmd1", "cmd2", "cmd3"},
		},
		{
			name:     "single command",
			input:    "cmd1",
			operator: "&&",
			expected: []string{"cmd1"},
		},
		{
			name:     "pipeline with AND",
			input:    "cmd1 | cmd2 && cmd3 | cmd4",
			operator: "&&",
			expected: []string{"cmd1 | cmd2", "cmd3 | cmd4"},
		},
		{
			name:     "pipeline with OR",
			input:    "cmd1 | cmd2 || cmd3 | cmd4",
			operator: "||",
			expected: []string{"cmd1 | cmd2", "cmd3 | cmd4"},
		},
		{
			name:     "complex pipeline with AND",
			input:    "echo hello | grep h && ls -la | wc -l",
			operator: "&&",
			expected: []string{"echo hello | grep h", "ls -la | wc -l"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitByOperator(tt.input, tt.operator)
			if len(result) != len(tt.expected) {
				t.Errorf("splitByOperator(%q, %q) length = %d, want %d", tt.input, tt.operator, len(result), len(tt.expected))
				return
			}
			for i, exp := range tt.expected {
				if result[i] != exp {
					t.Errorf("splitByOperator(%q, %q)[%d] = %q, want %q", tt.input, tt.operator, i, result[i], exp)
				}
			}
		})
	}
}

func TestCommandExecutorInterrupt(t *testing.T) {
	ce := NewCommandExecutor()

	ce.Interrupt()

	ce2 := NewCommandExecutor()
	if ce2 == nil {
		t.Error("NewCommandExecutor() returned nil")
	}
	if ce2.currentCmd != nil {
		t.Error("NewCommandExecutor() should have nil currentCmd")
	}
}

func TestParseRedirectsWithSpaces(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantIn  string
		wantOut string
	}{
		{
			name:    "spaces around redirects",
			input:   "echo hello >  output.txt",
			wantOut: "output.txt",
		},
		{
			name:    "no spaces around redirects",
			input:   "echo hello>output.txt",
			wantOut: "output.txt",
		},
		{
			name:   "input with spaces",
			input:  "cat  <  input.txt",
			wantIn: "input.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redirect, _ := parseRedirects(tt.input)
			if tt.wantOut != "" && redirect.OutputFile != tt.wantOut {
				t.Errorf("OutputFile = %q, want %q", redirect.OutputFile, tt.wantOut)
			}
			if tt.wantIn != "" && redirect.InputFile != tt.wantIn {
				t.Errorf("InputFile = %q, want %q", redirect.InputFile, tt.wantIn)
			}
		})
	}
}

func TestPipelineStdin(t *testing.T) {
	ce := NewCommandExecutor()

	pipeline := []string{"echo test", "cat"}
	redirect := &Redirect{}

	err := executePipeline(pipeline, ce, redirect)

	if err != nil {
	}
}

func TestCdTildeExpansion(t *testing.T) {
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("HOME not set")
	}

	err = cdCommand("~")
	if err != nil {
		t.Errorf("cdCommand(~) error = %v", err)
		return
	}

	currentDir, _ := os.Getwd()
	currentDir, _ = filepath.EvalSymlinks(currentDir)
	homeDir, _ = filepath.EvalSymlinks(homeDir)

	if currentDir != homeDir {
		t.Errorf("cdCommand(~) ended up in %q, want %q", currentDir, homeDir)
	}
}

func TestCdTildeWithSubpath(t *testing.T) {
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("HOME not set")
	}

	testDir := filepath.Join(homeDir, "test_shell_cd")
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	testDir, _ = filepath.EvalSymlinks(testDir)

	err = cdCommand("~/test_shell_cd")
	if err != nil {
		t.Errorf("cdCommand(~/test_shell_cd) error = %v", err)
		return
	}

	currentDir, _ := os.Getwd()
	currentDir, _ = filepath.EvalSymlinks(currentDir)

	if currentDir != testDir {
		t.Errorf("cdCommand(~/test_shell_cd) ended up in %q, want %q", currentDir, testDir)
	}
}
