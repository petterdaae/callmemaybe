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

	gen := language.NewAssemblyGenerator()
	gen.Start()
	err = ast.Generate(&gen)
	gen.End()
	if err != nil {
		return "", err
	}
	assembly := ""
	for i := range gen.Operations {
		assembly += "\t" + gen.Operations[i] + "\n"
	}

	for i := range gen.AllProcedures {
		proc := gen.AllProcedures[i]
		assembly += proc.Name + ":\n"
		for j := range proc.Operations {
			assembly += "\t" + proc.Operations[j] + "\n"
		}

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
