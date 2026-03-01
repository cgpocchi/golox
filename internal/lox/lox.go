package lox

import (
	"fmt"
	"os"
)

type ErrorTracker struct {
	HadError bool
}

func NewErrorTracker() *ErrorTracker {
	return &ErrorTracker{false}
}

func (e *ErrorTracker) Error(line int, message string) {
	Report(line, "", message)
	e.HadError = true
}

func Report(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
}
