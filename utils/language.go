package utils

import (
	"lang/language"
	"os/exec"
	"strings"
)

func Compile(program string) (string, error) {
	parser := language.NewParser(strings.NewReader(program))
	ast, err := parser.Parse()
	if err != nil {
		return "", err
	}

	out := language.AssemblyOutput{
		Operations: []string{},
		StackSize: 0,
		Identifiers: make(map[string]int),
	}
	out.Start()
	err = ast.Generate(&out)
	out.End()
	if err != nil {
		return "", err
	}
	assembly := ""
	for i := range out.Operations {
		assembly += out.Operations[i] + "\n"
	}
	return assembly, nil
}

func Assemble(file string) error {
	_, err := exec.Command("nasm", "-f", "elf64", file).CombinedOutput()
	return err
}

func Link(file string) error {
	_, err := exec.Command("gcc", "-no-pie", "-o", "out", file, "-lc").CombinedOutput()
	return err
}
