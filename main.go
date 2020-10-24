package main

import (
	"lang/grammar"
	"lang/parser"
	"strings"
)

func main() {
	program := "let a = 1 in let b = 2 in a + b * 3"
	parser := parser.New(strings.NewReader(program))
	exp, err := parser.ParseExp()
	if err != nil {
		println("Parser failed: " + err.Error())
		return
	}

	out := grammar.AssemblyOutput{
		Operations: []string{},
		StackSize: 0,
		Identifiers: make(map[string]int),
	}
	exp.Generate(&out)

	for i := range out.Operations {
		println(out.Operations[i])
	}
}
