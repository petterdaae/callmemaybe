package assemblyoutput

import (
	"fmt"
)

type AssemblyOutput struct {
	procedureStack       *ProcedureStack
	nameGeneratorCounter int
	EvaluatedProcedures  []*procedure
	MainOperations       []string
}

func NewAssemblyOutput() *AssemblyOutput {
	return &AssemblyOutput{
		procedureStack:       NewProcedureStack(),
		nameGeneratorCounter: 0,
		EvaluatedProcedures:  []*procedure{},
	}
}

func (ao *AssemblyOutput) Push(r string) {
	ao.addOperation(fmt.Sprintf("push %s", r))
}

func (ao *AssemblyOutput) Pop(r string) {
	ao.addOperation(fmt.Sprintf("pop %s", r))
}

func (ao *AssemblyOutput) Add(r1 string, r2 string) {
	ao.addOperation(fmt.Sprintf("add %s, %s", r1, r2))
}

func (ao *AssemblyOutput) Imul(r1 string, r2 string) {
	ao.addOperation(fmt.Sprintf("imul %s, %s", r1, r2))
}

func (ao *AssemblyOutput) Mov(r1 string, r2 string) {
	ao.addOperation(fmt.Sprintf("mov %s, %s", r1, r2))
}

func (ao *AssemblyOutput) Xor(r1 string, r2 string) {
	ao.addOperation(fmt.Sprintf("xor %s, %s", r1, r2))
}

func (ao *AssemblyOutput) CallPrintf() {
	ao.addOperation(fmt.Sprintf("call printf"))
}

func (ao *AssemblyOutput) Ret() {
	ao.addOperation(fmt.Sprintf("ret"))
}

func (ao *AssemblyOutput) Call(name string) {
	ao.addOperation(fmt.Sprintf("call %s", name))
}

func (ao *AssemblyOutput) Cmp(r1 string, r2 string) {
	ao.addOperation(fmt.Sprintf("cmp %s, %s", r1, r2))
}

func (ao *AssemblyOutput) Je(name string) {
	ao.addOperation(fmt.Sprintf("je %s", name))
}

func (ao *AssemblyOutput) Jne(name string) {
	ao.addOperation(fmt.Sprintf("jne %s", name))
}

func (ao *AssemblyOutput) Jg(name string) {
	ao.addOperation(fmt.Sprintf("jg %s", name))
}

func (ao *AssemblyOutput) Jl(name string) {
	ao.addOperation(fmt.Sprintf("jl %s", name))
}

func (ao *AssemblyOutput) Jle(name string) {
	ao.addOperation(fmt.Sprintf("jle %s", name))
}

func (ao *AssemblyOutput) Jge(name string) {
	ao.addOperation(fmt.Sprintf("jge %s", name))
}

func (ao *AssemblyOutput) Jmp(name string) {
	ao.addOperation(fmt.Sprintf("jmp %s", name))
}

func (ao *AssemblyOutput) Sub(r1, r2 string) {
	ao.addOperation(fmt.Sprintf("sub %s, %s", r1, r2))
}

func (ao *AssemblyOutput) Div(r string) {
	ao.addOperation(fmt.Sprintf("div %s", r))
}

func (ao *AssemblyOutput) NewSection(name string) {
	ao.addOperation(fmt.Sprintf("%s:", name))
}

func (ao *AssemblyOutput) addOperation(operation string) {
	procedure := ao.procedureStack.Peek()
	if procedure == nil {
		ao.MainOperations = append(ao.MainOperations, operation)
	} else {
		procedure.Operations = append(procedure.Operations, operation)
	}
}

func (ao *AssemblyOutput) PushProcedure(initialStackSize int, numberOfArguments int) string {
	name := ao.GenerateUniqueName()
	ao.procedureStack.Push(&procedure{
		Name:                              name,
		StackSizeBeforeFunctionGeneration: initialStackSize,
		NumberOfArgs:                      numberOfArguments,
	})
	return name
}

func (ao *AssemblyOutput) PopProcedure() {
	current := ao.procedureStack.Peek()
	ao.procedureStack.Pop()
	ao.EvaluatedProcedures = append(ao.EvaluatedProcedures, current)
}

func (ao *AssemblyOutput) GenerateUniqueName() string {
	ao.nameGeneratorCounter++
	return fmt.Sprintf("unique%d", ao.nameGeneratorCounter)
}

func (ao *AssemblyOutput) CurrentProcedure() *procedure {
	return ao.procedureStack.Peek()
}

func (ao *AssemblyOutput) Start() {
	ao.addOperation("extern printf")
	ao.addOperation("extern malloc")
	ao.addOperation("extern free")
	ao.addOperation("global main")
	ao.addOperation("section .date")
	ao.addOperation("digitNewlineFormat: db '%d', 10, 0")
	ao.addOperation("charNewlineFormat: db '%c', 10, 0")
	ao.addOperation("charFormat: db '%c', 0")
	ao.addOperation("section .text")
	ao.addOperation("main:")
	ao.addOperation("push rbx")
}

func (ao *AssemblyOutput) End(stackSize int) {
	ao.addOperation("pop rbx")
	for i := 0; i < stackSize; i++ {
		ao.addOperation("pop rbx")
	}
	ao.addOperation("mov rax, 0")
	ao.addOperation("ret")

	// Procedure for printing all registers in a list
	// RAX: list address, RBX: format
	ao.NewSection("printListWithFormat")
	ao.Mov(RDX, RAX)
	ao.Mov(RCX, "0")

	ao.Cmp(fmt.Sprintf("qword [%s]", RDX), "0")
	ao.Je("printListWithFormatLoopEnd")

	ao.NewSection("printListWithFormatLoopStart")
	ao.Add(RCX, "1")
	ao.Mov(RAX, fmt.Sprintf("[%s+%s*8]", RDX, RCX))
	ao.Call("printRegisterWithFormat")

	ao.Cmp(RCX, fmt.Sprintf("[%s]", RDX))
	ao.Je("printListWithFormatLoopEnd")
	ao.Jne("printListWithFormatLoopStart")

	ao.NewSection("printListWithFormatLoopEnd")
	ao.Ret()

	// Procedure for printing a register
	// RAX: register to print, RBX: format
	ao.NewSection("printRegisterWithFormat")
	ao.Push(RAX)
	ao.Push(RBX)
	ao.Push(RCX)
	ao.Push(RDX)
	ao.Mov(RDI, RBX)
	ao.Mov(RSI, RAX)
	ao.Xor(RAX, RAX)
	ao.CallPrintf()
	ao.Pop(RDX)
	ao.Pop(RCX)
	ao.Pop(RBX)
	ao.Pop(RAX)
	ao.Ret()
}
