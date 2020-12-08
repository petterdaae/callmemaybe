package assemblyoutput

import "fmt"

type AssemblyOutput struct {
	procedureStack ProcedureStack
	nameGeneratorCounter int
	evaluatedProcedures []*procedure
}

func (ao *AssemblyOutput) Push(r string) {
	fmt.Sprintf("push %s", r)
}

func (ao *AssemblyOutput) Pop(r string) {
	fmt.Sprintf("pop %s", r)
}

func (ao *AssemblyOutput) Add(r1 string, r2 string) {
	fmt.Sprintf("add %s, %s", r1, r2)
}

func (ao *AssemblyOutput) Imul(r1 string, r2 string) {
	fmt.Sprintf("imul %s, %s", r1, r2)
}

func (ao *AssemblyOutput) Mov(r1 string, r2 string) {
	fmt.Sprintf("mov %s, %s", r1, r2)
}

func (ao *AssemblyOutput) Xor(r1 string, r2 string) {
	fmt.Sprintf("xor %s, %s", r1, r2)
}

func (ao *AssemblyOutput) CallPrintf() {
	fmt.Sprintf("call printf")
}

func (ao *AssemblyOutput) Ret() {
	fmt.Sprintf("ret")
}

func (ao *AssemblyOutput) Call(name string) {
	fmt.Sprintf("call %s", name)
}

func (ao *AssemblyOutput) Cmp(r1 string, r2 string) {
	fmt.Sprintf("cmp %s, %s", r1, r2)
}

func (ao *AssemblyOutput) Je(name string) {
	fmt.Sprintf("je %s", name)
}

func (ao *AssemblyOutput) Jne(name string) {
	fmt.Sprintf("jne %s", name)
}

func (ao *AssemblyOutput) Jmp(name string) {
	fmt.Sprintf("jmp %s", name)
}

func (ao *AssemblyOutput) NewSection(name string) {
	fmt.Sprintf("%s:", name)
}

func (ao *AssemblyOutput) PushProcedure(numberOfArgs int, initialStackSize int) string {
	name := ao.GenerateUniqueName()
	ao.procedureStack.Push(&procedure{
		name:                              name,
		NumberOfArgs:                      numberOfArgs,
		StackSizeBeforeFunctionGeneration: initialStackSize,
	})
	return name
}

func (ao *AssemblyOutput) PopProcedure() {
	current := ao.procedureStack.Peek()
	ao.procedureStack.Pop()
	ao.evaluatedProcedures = append(ao.evaluatedProcedures, current)
}

func (ao *AssemblyOutput) GenerateUniqueName() string {
	ao.nameGeneratorCounter++
	return fmt.Sprintf("unique%d", ao.nameGeneratorCounter)
}

func (ao *AssemblyOutput) CurrentProcedure() *procedure {
	return ao.procedureStack.Peek()
}


