package main

import (
	"lang/parser"
	"strconv"
	"strings"
)

func main() {
	program := "(1 + 4) * 2 + 3 * 5"
	parser := parser.New(strings.NewReader(program))
	ast, err := parser.Parse()
	if err != nil {
		println("Parser failed: " + err.Error())
		return
	}
	value := ast.Evaluate()
	println("Result => " + strconv.Itoa(value))
}
