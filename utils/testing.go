package utils

import (
	"os"
	"testing"
)

func AssertProgramOutput(path string, output string, t *testing.T) {
	defer os.Remove("out")
	defer os.Remove("out.nasm")
	defer os.Remove("out.o")

	program, err := ReadFile(path)
	if err != nil {
		t.Errorf("failed to read file: %v", err)
		return
	}

	nasm, err := Compile(program)
	if err != nil {
		t.Errorf("failed to compile: %v", err)
		return
	}

	err = WriteFile("out.nasm", nasm)
	if err != nil {
		t.Errorf("failed to write file: %v", err)
		return
	}

	err = Assemble("out.nasm")
	if err != nil {
		t.Errorf("failed to assemble: %v", err)
		return
	}

	err = Link("./out.o")
	if err != nil {
		t.Errorf("failed to link: %v", err)
		return
	}

	stdout, err := RunExecutable("out")
	if err != nil {
		t.Errorf("failed to run executable: %v", err)
		return
	}

	if stdout != output {
		t.Errorf("got:\n%s\nexpected:\n%s\n", stdout, output)
	}
}

func AssertCompilerFails(path string, t *testing.T) {
	program, err := ReadFile(path)
	if err != nil {
		t.Errorf("failed to read file: %v", err)
		return
	}

	_, err = Compile(program)
	if err == nil {
		t.Errorf("did not fail to compile")
		return
	}
}
