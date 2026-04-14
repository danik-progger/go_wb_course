package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)

type CommandExecutor struct {
	currentCmd *exec.Cmd
}

func NewCommandExecutor() *CommandExecutor {
	return &CommandExecutor{}
}

func (ce *CommandExecutor) Interrupt() {
	if ce.currentCmd != nil && ce.currentCmd.Process != nil {
		ce.currentCmd.Process.Signal(syscall.SIGINT)
	}
}

func runShell(ce *CommandExecutor) {
	reader := bufio.NewReader(os.Stdin)

	for {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error getting current directory:", err)
			return
		}
		fmt.Printf("%s $> ", pwd)

		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("exit")
				return
			}
			fmt.Fprintln(os.Stderr, "error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		if err := executeInput(input, ce); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func splitByOperator(input, operator string) []string {
	var result []string
	var current strings.Builder

	rawParts := strings.Split(input, operator)

	for idx, part := range rawParts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if idx == 0 {
			current.WriteString(part)
		} else {
			if strings.Contains(current.String(), "|") {
				result = append(result, current.String())
				current.Reset()
				current.WriteString(part)
			} else if strings.Contains(part, "|") {
				if current.Len() > 0 {
					result = append(result, current.String())
					current.Reset()
				}
				current.WriteString(part)
			} else {
				if current.Len() > 0 {
					result = append(result, current.String())
					current.Reset()
				}
				current.WriteString(part)
			}
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

func executeInput(input string, ce *CommandExecutor) error {
	input = expandEnvVars(input)

	if strings.Contains(input, "&&") {
		parts := splitByOperator(input, "&&")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			if err := executeSingleInput(part, ce); err != nil {
				return err
			}
		}
		return nil
	}

	if strings.Contains(input, "||") {
		parts := splitByOperator(input, "||")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			if err := executeSingleInput(part, ce); err != nil {
				continue
			}
			break
		}
		return nil
	}

	return executeSingleInput(input, ce)
}

func executeSingleInput(input string, ce *CommandExecutor) error {
	redirects, cleanRedirects := parseRedirects(input)
	input = cleanRedirects

	pipeline := strings.Split(input, "|")

	if len(pipeline) == 1 {
		args := strings.Fields(pipeline[0])
		if len(args) == 0 {
			return nil
		}

		switch args[0] {
		case "cd":
			if len(args) < 2 {
				return fmt.Errorf("cd: missing argument")
			}
			return cdCommand(args[1])
		case "exit":
			fmt.Println("exit")
			os.Exit(0)
		}
	}

	return executePipeline(pipeline, ce, redirects)
}

func cdCommand(path string) error {
	if path == "~" || strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cd: could not get home directory: %v", err)
		}
		if path == "~" {
			path = homeDir
		} else {
			path = filepath.Join(homeDir, path[2:])
		}
	}
	return os.Chdir(path)
}

func expandEnvVars(input string) string {
	re := regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)
	input = re.ReplaceAllStringFunc(input, func(match string) string {
		varName := match[2 : len(match)-1]
		return os.Getenv(varName)
	})

	re2 := regexp.MustCompile(`\$([A-Za-z_][A-Za-z0-9_]*)`)
	input = re2.ReplaceAllStringFunc(input, func(match string) string {
		varName := match[1:]
		return os.Getenv(varName)
	})

	return input
}

type Redirect struct {
	InputFile  string
	OutputFile string
	Append     bool
}

func parseRedirects(input string) (*Redirect, string) {
	redirect := &Redirect{}

	outputRe := regexp.MustCompile(`(>>?)\s*([^\s]+)`)
	matches := outputRe.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		if match[1] == ">>" {
			redirect.Append = true
		}
		redirect.OutputFile = match[2]
	}
	input = outputRe.ReplaceAllString(input, "")

	inputRe := regexp.MustCompile(`<\s*([^\s]+)`)
	matches = inputRe.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		redirect.InputFile = match[1]
	}
	input = inputRe.ReplaceAllString(input, "")

	return redirect, strings.TrimSpace(input)
}

func executePipeline(pipeline []string, ce *CommandExecutor, redirect *Redirect) error {
	cmds := make([]*exec.Cmd, len(pipeline))
	var err error

	for i, cmdStr := range pipeline {
		args := strings.Fields(cmdStr)
		if len(args) == 0 {
			return fmt.Errorf("pipeline: empty command")
		}

		switch args[0] {
		case "pwd":
			cmds[i] = exec.Command("pwd")
		case "echo":
			cmds[i] = exec.Command("echo", args[1:]...)
		case "kill":
			cmds[i] = exec.Command("kill", args[1:]...)
		case "ps":
			cmds[i] = exec.Command("ps", args[1:]...)
		default:
			cmds[i] = exec.Command(args[0], args[1:]...)
		}
	}

	for i := 0; i < len(cmds)-1; i++ {
		cmds[i+1].Stdin, err = cmds[i].StdoutPipe()
		if err != nil {
			return err
		}
	}

	if len(cmds) > 0 {
		lastCmd := cmds[len(cmds)-1]

		if redirect != nil && redirect.OutputFile != "" {
			var outFile *os.File
			if redirect.Append {
				outFile, err = os.OpenFile(redirect.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			} else {
				outFile, err = os.Create(redirect.OutputFile)
			}
			if err != nil {
				return fmt.Errorf("redirect: cannot open output file: %v", err)
			}
			defer outFile.Close()
			lastCmd.Stdout = outFile
		} else {
			lastCmd.Stdout = os.Stdout
		}

		if redirect != nil && redirect.InputFile != "" {
			inFile, err := os.Open(redirect.InputFile)
			if err != nil {
				return fmt.Errorf("redirect: cannot open input file: %v", err)
			}
			defer inFile.Close()
			cmds[0].Stdin = inFile
		} else {
			cmds[0].Stdin = os.Stdin
		}

		lastCmd.Stderr = os.Stderr
	}

	for i := len(cmds) - 1; i >= 0; i-- {
		ce.currentCmd = cmds[i]
		if err := cmds[i].Start(); err != nil {
			return err
		}
	}

	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			os.Exit(0)
		}
	}

	ce.currentCmd = nil
	return nil
}

func main() {
	ce := NewCommandExecutor()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for range c {
			ce.Interrupt()
			fmt.Println()
		}
	}()

	runShell(ce)
}
