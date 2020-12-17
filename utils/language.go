package utils

import (
	"callmemaybe/language"
	"callmemaybe/language/assemblyoutput"
	"callmemaybe/language/memorymodel"
	"os/exec"
	"strings"
)

func Compile(program string) (string, error) {
	parser := language.NewParser(strings.NewReader(program))
	ast, err := parser.Parse()
	if err != nil {
		return "", err
	}

	ao := assemblyoutput.NewAssemblyOutput()
	mm := memorymodel.NewMemoryModel()
	ao.Start()
	err = ast.Generate(ao, mm)
	ao.End(mm.CurrentStackSize)

	if err != nil {
		return "", err
	}
	assembly := ""
	for i := range ao.MainOperations {
		assembly += ao.MainOperations[i] + "\n"
	}

	for i := range ao.EvaluatedProcedures {
		proc := ao.EvaluatedProcedures[i]
		assembly += proc.Name + ":\n"
		for j := range proc.Operations {
			assembly += proc.Operations[j] + "\n"
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
