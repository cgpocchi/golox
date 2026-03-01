package main

import (
	"bufio"
	"fmt"
	"golox/internal/lox"
	"golox/internal/scanner"
	"os"
)

func runFile(path string, errTracker *lox.ErrorTracker) {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("error %v reading file %s", err, path))
	}
	run(string(content), errTracker)
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

		run(line, errTracker)
		errTracker.HadError = false
	}
}

func run(source string, errTracker *lox.ErrorTracker) {
	scanner := scanner.NewScanner(source, errTracker)
	tokens := scanner.ScanTokens()

	// print tokens
	for _, tok := range tokens {
		fmt.Println(tok)
	}
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
