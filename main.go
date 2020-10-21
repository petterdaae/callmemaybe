package main

import (
	"lang/grammar"
	"lang/parser"
	"strings"
)

func main() {
	program := "println 1 println 1 + 1 println 3 * 1"
	parser := parser.New(strings.NewReader(program))
	ast, err := parser.Parse()
	if err != nil {
		println("Parser failed: " + err.Error())
		return
	}
	ast.Execute(grammar.NewContext())
}
