package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

// runShell is the main function for the shell.
func runShell() {
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

		if err := executeInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

// executeInput parses and executes the user's input.
func executeInput(input string) error {

	pipeline := strings.Split(input, "|")

	// If there is no pipe, check for built-in commands that modify the shell's state.
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
			return os.Chdir(args[1])
		case "exit":
			// This will be handled by the EOF check in the main loop, but we can also have an explicit command.
			// To exit, we can return a special error or have a flag.
			// For now, let's just exit cleanly.
			fmt.Println("exit")
			os.Exit(0)
		}
	}

	// Execute one or more commands in a pipeline.
	return executePipeline(pipeline)
}

// executePipeline handles the execution of a pipeline of commands.
func executePipeline(pipeline []string) error {
	cmds := make([]*exec.Cmd, len(pipeline))
	var err error

	for i, cmdStr := range pipeline {
		args := strings.Fields(cmdStr)
		if len(args) == 0 {
			return fmt.Errorf("pipeline: empty command")
		}

		// Check for non-state-modifying built-ins.
		switch args[0] {
		case "pwd":
			// We can't run pwd as a simple built-in inside a pipeline easily.
			// For simplicity, we will execute it as an external command if it's not the only command.
			// A real shell would handle this more gracefully.
			cmds[i] = exec.Command("pwd")
		case "echo":
			// Similar to pwd, let's use the external echo for pipelines.
			cmds[i] = exec.Command("echo", args[1:]...)
		case "kill":
			// And kill as well.
			cmds[i] = exec.Command("kill", args[1:]...)
		case "ps":
			cmds[i] = exec.Command("ps", args[1:]...)
		default:
			cmds[i] = exec.Command(args[0], args[1:]...)
		}
	}

	// Wire up the pipeline.
	for i := 0; i < len(cmds)-1; i++ {
		cmds[i+1].Stdin, err = cmds[i].StdoutPipe()
		if err != nil {
			return err
		}
	}

	// Set the final output and error streams.
	if len(cmds) > 0 {
		cmds[len(cmds)-1].Stdout = os.Stdout
		cmds[len(cmds)-1].Stderr = os.Stderr
	}

	// Start commands in reverse order.
	for i := len(cmds) - 1; i >= 0; i-- {
		if err := cmds[i].Start(); err != nil {
			return err
		}
	}

	// Wait for commands to finish.
	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			// We don't return the error here because a non-zero exit code is not a shell error.
		}
	}

	return nil
}

func main() {
	// Set up signal handling for Ctrl+C.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// When Ctrl+C is pressed, a newline is printed to get a fresh prompt.
		fmt.Println()
	}()

	runShell()
}
