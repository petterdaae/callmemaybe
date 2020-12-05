package language

import (
	"fmt"
	"strconv"
)

const (
	rax = "rax"
	rbx = "rbx"
	rcx = "rcx"
	rdx = "rdx"
	rsp = "rsp" // stack pointer
	add = "add"
	sub = "sub"
)

type AssemblyGenerator struct {
	contexts             []Context
	procedureStack       []*AssemblyProcedure
	stackSize            int
	procedureNameCounter int
	anonymousProcedures  []*AssemblyProcedure
	Operations           []string
	AllProcedures        []*AssemblyProcedure
}

type AssemblyProcedure struct {
	Name                     string
	Operations               []string
	stackSizeWhenInitialized int
}

type Context struct {
	fields     map[string]int
	procedures map[string]string
}

func NewAssemblyGenerator() AssemblyGenerator {
	fields := make(map[string]int)
	procedures := make(map[string]string)
	return AssemblyGenerator{
		contexts: []Context{{
			fields:     fields,
			procedures: procedures,
		}},
		procedureStack:       []*AssemblyProcedure{},
		stackSize:            0,
		procedureNameCounter: 0,
		Operations:           []string{},
	}
}

func (gen *AssemblyGenerator) pushAnonymousProcedure(proc *AssemblyProcedure) {
	gen.anonymousProcedures = append(gen.anonymousProcedures, proc)
}

func (gen *AssemblyGenerator) popAnonymousProcedure() *AssemblyProcedure {
	n := len(gen.anonymousProcedures)
	temp := gen.anonymousProcedures[n-1]
	gen.anonymousProcedures = gen.anonymousProcedures[:n-1]
	return temp
}

func (gen *AssemblyGenerator) peekContext() Context {
	n := len(gen.contexts)
	return gen.contexts[n-1]
}

func (gen *AssemblyGenerator) get(field string) (ExpKind, string, error) {
	// Get the context
	context := gen.peekContext()

	// Get the stack time of the field
	stack, ok := context.fields[field]

	if ok {
		// Find the difference
		diff := (gen.stackSize - stack) * 8

		if diff < 0 {
			return InvalidExpKind, "", fmt.Errorf("the stack reference was negative: %d", diff)
		}

		// Special case if we are inside a procedure
		if len(gen.procedureStack) != 0 {
			procedure := gen.peekProcedure()
			if stack > procedure.stackSizeWhenInitialized {
				return StackExp, fmt.Sprintf("[%s+%d]", rsp, diff), nil
			}
			return StackExp, fmt.Sprintf("[%s+%s+%d+%d]", rsp, rcx, diff, 8), nil
		}

		return StackExp, fmt.Sprintf("[%s+%d]", rsp, diff), nil
	}
	proc, ok := context.procedures[field]
	if ok {
		return ProcExp, proc, nil
	}
	return InvalidExpKind, "", fmt.Errorf("field not available in current context: %s", field)
}

func (gen *AssemblyGenerator) pushToStack(field string) {
	gen.peekContext().fields[field] = gen.stackSize
}

func (gen *AssemblyGenerator) pushContext() {
	fieldsCopy := make(map[string]int)
	// size := len(gen.peekContext().fields)
	for k, _ := range gen.peekContext().fields {
		_, address, _ := gen.get(k)
		gen.move(rax, address)
		gen.push(rax)
		fieldsCopy[k] = gen.stackSize
	}
	proceduresCopy := make(map[string]string)
	for k, v := range gen.peekContext().procedures {
		proceduresCopy[k] = v
	}
	gen.contexts = append(gen.contexts, Context{fields: fieldsCopy, procedures: proceduresCopy})
}

func (gen *AssemblyGenerator) popContext() {
	n := len(gen.contexts)
	gen.contexts = gen.contexts[:n-1]
}

func (gen *AssemblyGenerator) nameLastAnonymousProc(name string) {
	proc := gen.popAnonymousProcedure()
	gen.peekContext().procedures[name] = proc.Name
}

func (gen *AssemblyGenerator) pushProcedure() {
	gen.procedureNameCounter++
	newProcedure := AssemblyProcedure{
		Name: fmt.Sprintf("proc%d", gen.procedureNameCounter),
		stackSizeWhenInitialized: gen.stackSize,
	}
	gen.procedureStack = append(gen.procedureStack, &newProcedure)
	gen.move(rdx, strconv.Itoa(gen.stackSize))
	gen.sub(rcx, rdx)
	gen.mult(rcx, "8")
}

func (gen *AssemblyGenerator) peekProcedure() *AssemblyProcedure {
	n := len(gen.procedureStack)
	return gen.procedureStack[n-1]
}

func (gen *AssemblyGenerator) popProcedure() {
	n := len(gen.procedureStack)
	last := gen.procedureStack[n-1]

	diff := gen.stackSize - last.stackSizeWhenInitialized
	for i := 0; i < diff; i++ {
		gen.pop(rax)
	}

	gen.ret()

	gen.anonymousProcedures = append(gen.anonymousProcedures, last)
	gen.AllProcedures = append(gen.AllProcedures, last)
	gen.procedureStack = gen.procedureStack[:n-1]
}

func (gen *AssemblyGenerator) Start() {
	gen.Operations = append(gen.Operations, "extern printf")
	gen.Operations = append(gen.Operations, "global main")
	gen.Operations = append(gen.Operations, "section .data")
	gen.Operations = append(gen.Operations, "format: db '%d', 10, 0")
	gen.Operations = append(gen.Operations, "section .text")
	gen.Operations = append(gen.Operations, "main:")
	gen.Operations = append(gen.Operations, "push rbx") // stack pointer might not be initialized without doing this?
}

func (gen *AssemblyGenerator) End() {
	gen.Operations = append(gen.Operations, "pop rbx")
	for i := 0; i < gen.stackSize; i++ {
		gen.Operations = append(gen.Operations, "pop rbx") // stack should be empty at end of program? (fixed segmentation fault)
	}
	gen.Operations = append(gen.Operations, "mov rax, 0")
	gen.Operations = append(gen.Operations, "ret")
}

func (gen *AssemblyGenerator) addOperation(operation string) {
	if len(gen.procedureStack) > 0 {
		procedure := gen.peekProcedure()
		procedure.Operations = append(procedure.Operations, operation)
	} else {
		gen.Operations = append(gen.Operations, operation)
	}
}

func (gen *AssemblyGenerator) move(destination string, source string) {
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

func (gen *AssemblyGenerator) pop(destination string) {
	line := fmt.Sprintf("pop %s", destination)
	gen.stackSize--
	gen.addOperation(line)
}

func (gen *AssemblyGenerator) sub(destination string, value string) {
	line := fmt.Sprintf("%s %s, %s", sub, destination, value)
	gen.addOperation(line)
}

func (gen *AssemblyGenerator) ret() {
	gen.addOperation("ret")
}

func (gen *AssemblyGenerator) call(name string) error {
	kind, actualName, err := gen.get(name)
	if kind != ProcExp || err != nil {
		return fmt.Errorf("failed to call procedure with name: %s", name)
	}

	var procedure *AssemblyProcedure
	for i := range gen.AllProcedures {
		p := gen.AllProcedures[i]
		if p.Name == actualName {
			procedure = p
		}
	}

	if procedure == nil {
		return fmt.Errorf("failed to find procedure with name: %s", actualName)
	}

	gen.move(rcx, strconv.Itoa(gen.stackSize))


	gen.addOperation(fmt.Sprintf("call %s", actualName))
	return nil
}

func (gen *AssemblyGenerator) println(value string) {
	gen.addOperation(fmt.Sprintf("mov rdi, format"))
	gen.addOperation(fmt.Sprintf("mov rsi, %s", value))
	gen.addOperation(fmt.Sprintf("xor rax, rax"))
	gen.addOperation(fmt.Sprintf("call printf"))
}
