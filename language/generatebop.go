package language

import (
	"callmemaybe/language/assemblyoutput"
	"callmemaybe/language/common"
	"callmemaybe/language/memorymodel"
	"fmt"
)

func (exp ExpGreater) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerateResult {
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
	isValidKind := func(gr GenerateResult) bool {
		return gr.Kind.IsComparable()
	}
	return HelpGenerateStackBop(ao, mm, exp, "greater", operation, KindBool, isValidKind)
}

func (exp ExpLess) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerateResult {
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
	isValidKind := func(gr GenerateResult) bool {
		return gr.Kind.IsComparable()
	}
	return HelpGenerateStackBop(ao, mm, exp, "less", operation, KindBool, isValidKind)
}

func (exp ExpEquals) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerateResult {
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
	isValidKind := func(gr GenerateResult) bool {
		return gr.Kind.IsComparable()
	}
	return HelpGenerateStackBop(ao, mm, exp, "equals", operation, KindBool, isValidKind)
}

func (exp ExpPlus) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerateResult {
	operation := func(ao *assemblyoutput.AssemblyOutput) {
		ao.Add(RAX, RBX)
	}
	isValidKind := func(gr GenerateResult) bool {
		return gr.Kind.IsAlgebraic()
	}
	return HelpGenerateStackBop(ao, mm, exp, "plus", operation, KindNumber, isValidKind)
}

func (exp ExpMultiply) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerateResult {
	operation := func(ao *assemblyoutput.AssemblyOutput) {
		ao.Imul(RAX, RBX)
	}
	isValidKind := func(gr GenerateResult) bool {
		return gr.Kind.IsAlgebraic()
	}
	return HelpGenerateStackBop(ao, mm, exp, "multiply", operation, KindNumber, isValidKind)
}

func (exp ExpMinus) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerateResult {
	operation := func(ao *assemblyoutput.AssemblyOutput) {
		ao.Sub(RBX, RAX)
		ao.Mov(RAX, RBX)
	}
	isValidKind := func(gr GenerateResult) bool {
		return gr.Kind.IsAlgebraic()
	}
	return HelpGenerateStackBop(ao, mm, exp, "multiply", operation, KindNumber, isValidKind)
}

func (exp ExpDivide) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerateResult {
	operation := func(ao *assemblyoutput.AssemblyOutput) {
		ao.Mov(RDX, "0")
		ao.Mov(RCX, RAX)
		ao.Mov(RAX, RBX)
		ao.Div(RCX)
	}
	isValidKind := func(gr GenerateResult) bool {
		return gr.Kind.IsAlgebraic()
	}
	return HelpGenerateStackBop(ao, mm, exp, "divide", operation, KindNumber, isValidKind)
}

func (exp ExpModulo) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerateResult {
	operation := func(ao *assemblyoutput.AssemblyOutput) {
		ao.Mov(RDX, "0")
		ao.Mov(RCX, RAX)
		ao.Mov(RAX, RBX)
		ao.Div(RCX)
		ao.Mov(RAX, RDX)
	}
	isValidKind := func(gr GenerateResult) bool {
		return gr.Kind.IsAlgebraic()
	}
	return HelpGenerateStackBop(ao, mm, exp, "modulo", operation, KindNumber, isValidKind)
}

func HelpGenerateStackBop(
	ao *assemblyoutput.AssemblyOutput,
	mm *memorymodel.MemoryModel,
	exp ExpBop,
	name string,
	operation func(ao *assemblyoutput.AssemblyOutput),
	kind common.ContextElementKind,
	isValidKind func(result GenerateResult) bool,
) GenerateResult {
	resultLeft := exp.LeftExp().Generate(ao, mm)
	if resultLeft.IsError(){
		return resultLeft.WrapError(fmt.Sprintf("failed to generate left expression of %s", name))
	}
	if !isValidKind(resultLeft) {
		return ErrorResult(fmt.Errorf("invalidl kind at left side of %s expression", name))
	}

	mm.CurrentStackSize++
	ao.Push(RAX)

	resultRight := exp.RightExp().Generate(ao, mm)
	if resultRight.IsError() {
		return resultRight.WrapError(fmt.Sprintf("failed to generate right expression of %s", name))
	}
	if !isValidKind(resultRight) {
		return ErrorResult(fmt.Errorf("%s only supports stack kinds", name))
	}

	if resultLeft.Kind != resultRight.Kind {
		return ErrorResult(fmt.Errorf("mismatching kinds in %s expression", name))
	}
	mm.CurrentStackSize--
	ao.Pop(RBX)
	operation(ao)
	return CustomResult(kind, KindInvalid, "", nil)
}
