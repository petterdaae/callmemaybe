package main

import (
	"lang/grammar"
	"lang/parser"
	"strings"
)

func main() {
	program := `
				a = 3
				println a
				b = 154 + let c = 2 in c + 3
				println a + b
				n = a * b
				println n
               `
	parser := parser.New(strings.NewReader(program))
	exp, err := parser.Parse()
	if err != nil {
		println("Parser failed: " + err.Error())
		return
	}

	out := grammar.AssemblyOutput{
		Operations: []string{},
		StackSize: 0,
		Identifiers: make(map[string]int),
	}
	out.Start()
	err = exp.Generate(&out)
	out.End()

	if err != nil {
		println(err.Error())
		return
	}

	for i := range out.Operations {
		println(out.Operations[i])
	}
}
