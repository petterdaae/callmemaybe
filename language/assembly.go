package language

import (
	"fmt"
	"lang/language/assembly"
)

const (
	rax = "rax"
	rbx = "rbx"
	rcx = "rcx"
	rdx = "rdx"
	rsp = "rsp"
	add = "add"
	sub = "sub"
)

type AssemblyGenerator struct {
	contexts             *assembly.Contexts
	stackSize            int
	Operations           []string
	AllProcedures        []*assembly.Procedure
	NamelessProcedures   []*assembly.Procedure
}

type Context struct {
	fields     map[string]int
	procedures map[string]string
}

func NewAssemblyGenerator() AssemblyGenerator {
	return AssemblyGenerator{
		contexts: assembly.NewContexts(),
		stackSize:            0,
		Operations:           []string{},
	}
}

func (gen *AssemblyGenerator) Start() {
	gen.Operations = append(gen.Operations, "extern printf")
	gen.Operations = append(gen.Operations, "global main")
	gen.Operations = append(gen.Operations, "section .data")
	gen.Operations = append(gen.Operations, "format: db '%d', 10, 0")
	gen.Operations = append(gen.Operations, "section .text")
	gen.Operations = append(gen.Operations, "main:")
	gen.Operations = append(gen.Operations, "push rbx")
}

func (gen *AssemblyGenerator) End() {
	gen.Operations = append(gen.Operations, "pop rbx")
	for i := 0; i < gen.stackSize; i++ {
		gen.Operations = append(gen.Operations, "pop rbx")
	}
	gen.Operations = append(gen.Operations, "mov rax, 0")
	gen.Operations = append(gen.Operations, "ret")
}

func (gen *AssemblyGenerator) addOperation(operation string) {
	procedure := gen.contexts.GetTopProcedure()
	if procedure != nil {
		procedure.Operations = append(procedure.Operations, operation)
	} else {
		gen.Operations = append(gen.Operations, operation)
	}
}

func (gen *AssemblyGenerator) mov(destination string, source string) {
	line := fmt.Sprintf("mov %s, %s", destination, source)
	gen.addOperation(line)
}

func (gen *AssemblyGenerator) add(destination string, value string) {
	line := fmt.Sprintf("%s %s, %s", add, destination, value)
	gen.addOperation(line)
}

func (gen *AssemblyGenerator) mult(destination string, value string) {
	line := fmt.Sprintf("imul %s, %s", destination, value)
	gen.addOperation(line)
}

func (gen *AssemblyGenerator) push(value string) {
	line := fmt.Sprintf("push %s", value)
	gen.stackSize++
	gen.addOperation(line)
}


func (gen *AssemblyGenerator) pushWithoutIncreasingStackSize(value string) {
	line := fmt.Sprintf("push %s", value)
	gen.addOperation(line)
}

func (gen *AssemblyGenerator) pop(destination string) {
	line := fmt.Sprintf("pop %s", destination)
	gen.stackSize--
	gen.addOperation(line)
}


func (gen *AssemblyGenerator) popWithoutDecreasingStackSize(destination string) {
	line := fmt.Sprintf("pop %s", destination)
	gen.addOperation(line)
}

func (gen *AssemblyGenerator) sub(destination string, value string) {
	line := fmt.Sprintf("%s %s, %s", sub, destination, value)
	gen.addOperation(line)
}

func (gen *AssemblyGenerator) ret() {
	gen.addOperation("ret")
}

func (gen *AssemblyGenerator) AddOperations(operations []string) {
	for _, operation := range operations {
		gen.addOperation(operation)
	}
}

func (gen *AssemblyGenerator) PushNamelessProcedure(procedure *assembly.Procedure) {
	gen.NamelessProcedures = append(gen.NamelessProcedures, procedure)
	gen.AllProcedures = append(gen.AllProcedures, procedure)
}

func (gen *AssemblyGenerator) NameNamelessProcedure(name string) {
	n := len(gen.NamelessProcedures)
	recent := gen.NamelessProcedures[n-1]
	gen.contexts.ProcInsert(name, recent.Name)
	gen.NamelessProcedures = gen.NamelessProcedures[:n-1]
}

func (gen *AssemblyGenerator) println(value string) {
	gen.addOperation(fmt.Sprintf("mov rdi, format"))
	gen.addOperation(fmt.Sprintf("mov rsi, %s", value))
	gen.addOperation(fmt.Sprintf("xor rax, rax"))
	gen.addOperation(fmt.Sprintf("call printf"))
}
