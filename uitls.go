package main

import (
	"io/ioutil"
	"lang/grammar"
	"lang/parser"
	"os"
	"strings"
)

func ReadFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func Compile(program string) (string, error) {
	println("    Parsing ...")
	parser := parser.New(strings.NewReader(program))
	ast, err := parser.Parse()
	if err != nil {
		println("❌ Parsing failed")
		return "", err
	}

	println("    Compiling ...")
	out := grammar.AssemblyOutput{
		Operations: []string{},
		StackSize: 0,
		Identifiers: make(map[string]int),
	}
	out.Start()
	err = ast.Generate(&out)
	out.End()
	if err != nil {
		println("❌ Compilation failed")
		return "", err
	}
	assembly := ""
	for i := range out.Operations {
		assembly += out.Operations[i] + "\n"
	}
	return assembly, nil
}

func WriteFile(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}
