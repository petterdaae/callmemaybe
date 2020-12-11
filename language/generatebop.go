package language

import (
	"fmt"
	"lang/language/assemblyoutput"
	"lang/language/memorymodel"
)

func (exp ExpGreater) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerationResult {
	operation := func(ao *assemblyoutput.AssemblyOutput) {
		ao.Cmp(RBX, RAX)
		greater := ao.GenerateUniqueName()
		lessThanOrEqual := ao.GenerateUniqueName()
		done := ao.GenerateUniqueName()
		ao.Jg(greater)
		ao.Jle(lessThanOrEqual)
		ao.NewSection(greater)
		ao.Mov(RAX, "1")
		ao.Jmp(done)
		ao.NewSection(lessThanOrEqual)
		ao.Mov(RAX, "0")
		ao.NewSection(done)
	}
	return HelpGenerateStackBop(ao, mm, exp, "greater", operation, KindBool)
}

func (exp ExpLess) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerationResult {
	operation := func(ao *assemblyoutput.AssemblyOutput) {
		ao.Cmp(RBX, RAX)
		less := ao.GenerateUniqueName()
		greaterThanOrEqual := ao.GenerateUniqueName()
		done := ao.GenerateUniqueName()
		ao.Jl(less)
		ao.Jge(greaterThanOrEqual)
		ao.NewSection(less)
		ao.Mov(RAX, "1")
		ao.Jmp(done)
		ao.NewSection(greaterThanOrEqual)
		ao.Mov(RAX, "0")
		ao.NewSection(done)
	}
	return HelpGenerateStackBop(ao, mm, exp, "less", operation, KindBool)
}

func (exp ExpEquals) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerationResult {
	operation := func(ao *assemblyoutput.AssemblyOutput) {
		ao.Cmp(RBX, RAX)
		equal := ao.GenerateUniqueName()
		notEqual := ao.GenerateUniqueName()
		done := ao.GenerateUniqueName()
		ao.Je(equal)
		ao.Jne(notEqual)
		ao.NewSection(equal)
		ao.Mov(RAX, "1")
		ao.Jmp(done)
		ao.NewSection(notEqual)
		ao.Mov(RAX, "0")
		ao.NewSection(done)
	}
	return HelpGenerateStackBop(ao, mm, exp, "equals", operation, KindBool)
}

func (exp ExpPlus) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerationResult {
	operation := func(ao *assemblyoutput.AssemblyOutput) {
		ao.Add(RAX, RBX)
	}
	return HelpGenerateStackBop(ao, mm, exp, "plus", operation, KindNumber)
}

func (exp ExpMultiply) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerationResult {
	operation := func(ao *assemblyoutput.AssemblyOutput) {
		ao.Imul(RAX, RBX)
	}
	return HelpGenerateStackBop(ao, mm, exp, "multiply", operation, KindNumber)
}

func (exp ExpMinus) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerationResult {
	operation := func(ao *assemblyoutput.AssemblyOutput) {
		ao.Sub(RBX, RAX)
		ao.Mov(RAX, RBX)
	}
	return HelpGenerateStackBop(ao, mm, exp, "multiply", operation, KindNumber)
}

func (exp ExpDivide) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerationResult {
	operation := func(ao *assemblyoutput.AssemblyOutput) {
		ao.Mov(RDX, "0")
		ao.Mov(RCX, RAX)
		ao.Mov(RAX, RBX)
		ao.Div(RCX)
	}
	return HelpGenerateStackBop(ao, mm, exp, "divide", operation, KindNumber)
}

func (exp ExpModulo) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerationResult {
	operation := func(ao *assemblyoutput.AssemblyOutput) {
		ao.Mov(RDX, "0")
		ao.Mov(RCX, RAX)
		ao.Mov(RAX, RBX)
		ao.Div(RCX)
		ao.Mov(RAX, RDX)
	}
	return HelpGenerateStackBop(ao, mm, exp, "modulo", operation, KindNumber)
}

func HelpGenerateStackBop(
	ao *assemblyoutput.AssemblyOutput,
	mm *memorymodel.MemoryModel,
	exp ExpBop,
	name string,
	operation func(ao *assemblyoutput.AssemblyOutput),
	kind memorymodel.ContextElementKind,
) GenerationResult {
	result := exp.LeftExp().Generate(ao, mm)
	if result.Error != nil {
		return ErrorKind(fmt.Errorf("failed to generate left expression of %s: %w", name, result.Error))
	}
	if !memorymodel.IsIntOrBool(kind) {
		return ErrorKind(fmt.Errorf("%s only supports stack kinds", name))
	}

	mm.CurrentStackSize++
	ao.Push(RAX)

	result = exp.RightExp().Generate(ao, mm)
	if result.Error != nil {
		return ErrorKind(fmt.Errorf("failed to generate right expression of %s: %w", name, result.Error))
	}
	if !memorymodel.IsIntOrBool(kind) {
		return ErrorKind(fmt.Errorf("%s only supports stack kinds", name))
	}

	mm.CurrentStackSize--
	ao.Pop(RBX)

	operation(ao)

	return NumberKind()
}
