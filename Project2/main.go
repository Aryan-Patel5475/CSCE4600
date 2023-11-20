package main

import (
	"bufio"
	"fmt"
	"github.com/Aryan-Patel5475/CSCE4600/Project2/builtins"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

func main() {
	exit := make(chan struct{}, 2) // buffer this so there's no deadlock.
	runLoop(os.Stdin, os.Stdout, os.Stderr, exit)
}

func runLoop(r io.Reader, w, errW io.Writer, exit chan struct{}) {
	var (
		input    string
		err      error
		readLoop = bufio.NewReader(r)
	)
	for {
		select {
		case <-exit:
			_, _ = fmt.Fprintln(w, "exiting gracefully...")
			return
		default:
			if err := printPrompt(w); err != nil {
				_, _ = fmt.Fprintln(errW, err)
				continue
			}
			if input, err = readLoop.ReadString('\n'); err != nil {
				_, _ = fmt.Fprintln(errW, err)
				continue
			}
			if err = handleInput(w, input, exit); err != nil {
				_, _ = fmt.Fprintln(errW, err)
			}
		}
	}
}

func printPrompt(w io.Writer) error {
	// Get current user.
	// Don't prematurely memoize this because it might change due to `su`?
	u, err := user.Current()
	if err != nil {
		return err
	}
	// Get current working directory.
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// /home/User [Username] $
	_, err = fmt.Fprintf(w, "%v [%v] $ ", wd, u.Username)

	return err
}

func handleInput(w io.Writer, input string, exit chan<- struct{}) error {
	// Remove trailing spaces.
	input = strings.TrimSpace(input)

	// Split the input separate the command name and the command arguments.
	args := strings.Split(input, " ")
	name, args := args[0], args[1:]

	// Check for built-in commands.
	// New builtin commands should be added here. Eventually this should be refactored to its own func.
	return handleBuiltinCommand(w, name, args, exit)
}

func handleBuiltinCommand(w io.Writer, name string, args []string, exit chan<- struct{}) error {
	switch name {
	case "cd":
		return builtins.ChangeDirectory(args...)
	case "env":
		return builtins.EnvironmentVariables(w, args...)
	case "whoami":
		u, err := user.Current()
		if err != nil {
			return err
		}
		fmt.Fprintln(w, u.Username)
		return nil
	case "pwd":
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		fmt.Fprintln(w, wd)
		return nil
	case "echo":
		fmt.Fprintln(w, strings.Join(args, " "))
		return nil
	case "mkdir":
		if len(args) == 0 {
			return fmt.Errorf("mkdir: missing operand")
		}
		for _, dir := range args {
			if err := os.Mkdir(dir, 0755); err != nil {
				return err
			}
		}
		return nil
	case "rm":
		if len(args) == 0 {
			return fmt.Errorf("rm: missing operand")
		}
		for _, file := range args {
			if err := os.Remove(file); err != nil {
				return err
			}
		}
		return nil
	case "exit":
		exit <- struct{}{}
		return nil
	}

	return executeCommand(name, args...)
}

func executeCommand(name string, arg ...string) error {
	// Otherwise prep the command
	cmd := exec.Command(name, arg...)

	// Set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Execute the command and return the error.
	return cmd.Run()
}
