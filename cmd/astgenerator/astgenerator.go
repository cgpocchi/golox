package main

import (
	"fmt"
	"iter"
	"maps"
	"os"
	"strings"
)

// Write a go file that defines all the expression in our AST.
func defineAst(outputDir, baseName string, types map[string][]string) {
	path := fmt.Sprintf("%s/%s.go", outputDir, baseName)
	f, err := os.Create(path)
	if err != nil {
		panic(fmt.Sprintf("Failed to create file at path %s", path))
	}
	defer f.Close()

	// setup package and imports
	writeToFile(f, fmt.Sprintf("package %s\n", baseName))
	writeToFile(f, "\n")
	writeToFile(f, "import \"golox/internal/token\"\n")
	writeToFile(f, "\n")

	// define visitor interface
	defineVisitor(f, maps.Keys(types))

	// add Expr interface
	writeToFile(f, "type Expr interface {\n")
	writeToFile(f, "    Accept(visitor Visitor[any]) any\n")
	writeToFile(f, "}\n\n")

	// define structs for each CFG rule
	for structName, fields := range types {
		writeToFile(f, "\n")
		defineType(f, structName, fields)
	}
}

func defineVisitor(f *os.File, types iter.Seq[string]) {
	writeToFile(f, "type Visitor[T any] interface {\n")
	for exprType := range types {
		writeToFile(f, fmt.Sprintf("    Visit%s(expr *%s) T\n", exprType, exprType))
	}
	writeToFile(f, "}\n\n")
}

// Create a Type struct in the given file.
func defineType(f *os.File, structName string, fieldList []string) {
	// create struct
	writeToFile(f, fmt.Sprintf("type %s struct {\n", structName))
	var sb strings.Builder
	for i, field := range fieldList {
		writeToFile(f, fmt.Sprintf("    %s\n", field))
		if i != len(fieldList)-1 {
			sb.WriteString(fmt.Sprintf("%s, ", field))
		} else {
			sb.WriteString(field)
		}
	}
	writeToFile(f, "}\n\n")

	// create constructor method
	writeToFile(f, fmt.Sprintf("func New%s(%s) (* %s) {\n", structName, sb.String(), structName))
	writeToFile(f, fmt.Sprintf("    return &%s{\n", structName))
	for _, field := range fieldList {
		param := strings.Split(field, " ")[0]
		writeToFile(f, fmt.Sprintf("        %s: %s,\n", param, param))
	}
	writeToFile(f, "    }\n")
	writeToFile(f, "}\n\n")

	// Implement visitor accept method
	writeToFile(f, fmt.Sprintf("func (e *%s) Accept(visitor Visitor[any]) any {\n", structName))
	writeToFile(f, fmt.Sprintf("    return visitor.Visit%s(e)\n", structName))
	writeToFile(f, "}\n\n")
}

// Write the given string to the given file.
func writeToFile(file *os.File, content string) {
	_, err := fmt.Fprint(file, content)
	if err != nil {
		panic("Failed to write AST content to file")
	}
}

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: generate_ast <output directory>")
		os.Exit(64)
	}
	outputDir := args[1]
	defineAst(outputDir, "expression", map[string][]string{
		"Binary":   []string{"left Expr", "operator token.Token", "right Expr"},
		"Grouping": []string{"expression Expr"},
		"Literal":  []string{"value any"},
		"Unary":    []string{"operator token.Token", "right Expr"},
	})
}
