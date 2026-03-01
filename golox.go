package main

import (
	"bufio"
	"fmt"
	"golox/internal/lox"
	"os"
)

func runFile(path string, errTracker *lox.ErrorTracker) {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("error %v reading file %s", err, path))
	}
	run(string(content))
	if errTracker.HadError {
		os.Exit(65)
	}
}

func runPrompt(errTracker *lox.ErrorTracker) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			panic("Failed to read prompt")
		}
		if line == "" {
			break
		}

		run(line)
		errTracker.HadError = false
	}
}

// TODO: Implement Me!
func run(source string) {
	fmt.Println("Echoing command: " + source)
}

func main() {
	args := os.Args
	errTracker := lox.NewErrorTracker()
	switch n := len(args); n {
	case 1:
		runPrompt(errTracker)
	case 2:
		runFile(args[1], errTracker)
	default:
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	}
}
