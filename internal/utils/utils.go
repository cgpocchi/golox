// Package for general purpose utilities.
package utils

// Ternary operator implementation. Equivalent to cond ? a : b.
func TernaryOp[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
