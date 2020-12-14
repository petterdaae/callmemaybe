package language

import (
	"callmemaybe/language/assemblyoutput"
	"callmemaybe/language/memorymodel"
	"callmemaybe/language/typesystem"
	"fmt"
)

func (exp ExpGreater) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (typesystem.Type, error) {
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
	isValidKind := func(gr typesystem.Type) bool {
		return gr.IsComparable()
	}
	return HelpGenerateStackBop(ao, mm, exp, "greater", operation, typesystem.NewBool(), isValidKind)
}

func (exp ExpLess) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (typesystem.Type, error) {
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
	isValidKind := func(gr typesystem.Type) bool {
		return gr.IsComparable()
	}
	return HelpGenerateStackBop(ao, mm, exp, "less", operation, typesystem.NewBool(), isValidKind)
}

func (exp ExpEquals) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (typesystem.Type, error) {
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
	isValidKind := func(gr typesystem.Type) bool {
		return gr.IsComparable()
	}
	return HelpGenerateStackBop(ao, mm, exp, "equals", operation, typesystem.NewBool(), isValidKind)
}

func (exp ExpNotEquals) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (typesystem.Type, error) {
	operation := func(ao *assemblyoutput.AssemblyOutput) {
		ao.Cmp(RBX, RAX)
		equal := ao.GenerateUniqueName()
		notEqual := ao.GenerateUniqueName()
		done := ao.GenerateUniqueName()
		ao.Je(equal)
		ao.Jne(notEqual)
		ao.NewSection(equal)
		ao.Mov(RAX, "0")
		ao.Jmp(done)
		ao.NewSection(notEqual)
		ao.Mov(RAX, "1")
		ao.NewSection(done)
	}
	isValidKind := func(gr typesystem.Type) bool {
		return gr.IsComparable()
	}
	return HelpGenerateStackBop(ao, mm, exp, "equals", operation, typesystem.NewBool(), isValidKind)
}

func (exp ExpPlus) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (typesystem.Type, error) {
	operation := func(ao *assemblyoutput.AssemblyOutput) {
		ao.Add(RAX, RBX)
	}
	isValidKind := func(gr typesystem.Type) bool {
		return gr.IsAlgebraic()
	}
	return HelpGenerateStackBop(ao, mm, exp, "plus", operation, typesystem.NewInt(), isValidKind)
}

func (exp ExpMultiply) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (typesystem.Type, error) {
	operation := func(ao *assemblyoutput.AssemblyOutput) {
		ao.Imul(RAX, RBX)
	}
	isValidKind := func(gr typesystem.Type) bool {
		return gr.IsAlgebraic()
	}
	return HelpGenerateStackBop(ao, mm, exp, "multiply", operation, typesystem.NewInt(), isValidKind)
}

func (exp ExpMinus) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (typesystem.Type, error) {
	operation := func(ao *assemblyoutput.AssemblyOutput) {
		ao.Sub(RBX, RAX)
		ao.Mov(RAX, RBX)
	}
	isValidKind := func(gr typesystem.Type) bool {
		return gr.IsAlgebraic()
	}
	return HelpGenerateStackBop(ao, mm, exp, "multiply", operation, typesystem.NewInt(), isValidKind)
}

func (exp ExpDivide) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (typesystem.Type, error) {
	operation := func(ao *assemblyoutput.AssemblyOutput) {
		ao.Mov(RDX, "0")
		ao.Mov(RCX, RAX)
		ao.Mov(RAX, RBX)
		ao.Div(RCX)
	}
	isValidKind := func(gr typesystem.Type) bool {
		return gr.IsAlgebraic()
	}
	return HelpGenerateStackBop(ao, mm, exp, "divide", operation, typesystem.NewInt(), isValidKind)
}

func (exp ExpModulo) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (typesystem.Type, error) {
	operation := func(ao *assemblyoutput.AssemblyOutput) {
		ao.Mov(RDX, "0")
		ao.Mov(RCX, RAX)
		ao.Mov(RAX, RBX)
		ao.Div(RCX)
		ao.Mov(RAX, RDX)
	}
	isValidKind := func(gr typesystem.Type) bool {
		return gr.IsAlgebraic()
	}
	return HelpGenerateStackBop(ao, mm, exp, "modulo", operation, typesystem.NewInt(), isValidKind)
}

func HelpGenerateStackBop(
	ao *assemblyoutput.AssemblyOutput,
	mm *memorymodel.MemoryModel,
	exp ExpBop,
	name string,
	operation func(ao *assemblyoutput.AssemblyOutput),
	kind typesystem.Type,
	isValidKind func(kind typesystem.Type) bool,
) (typesystem.Type, error) {
	kindLeft, err := exp.LeftExp().Generate(ao, mm)
	if err != nil {
		return typesystem.NewInvalid(), fmt.Errorf("failed to generate left expression of %s", name)
	}
	if !isValidKind(kindLeft) {
		return typesystem.NewInvalid(), fmt.Errorf("invalidl type at left side of %s expression", name)
	}

	mm.CurrentStackSize++
	ao.Push(RAX)

	kindRight, err := exp.RightExp().Generate(ao, mm)
	if err != nil {
		return typesystem.NewInvalid(), fmt.Errorf("failed to generate right expression of %s", name)
	}
	if !isValidKind(kindRight) {
		return typesystem.NewInvalid(), fmt.Errorf("%s only supports stack kinds", name)
	}

	if !kindLeft.Equals(kindRight) {
		return typesystem.NewInvalid(), fmt.Errorf("mismatching kinds in %s expression", name)
	}
	mm.CurrentStackSize--
	ao.Pop(RBX)
	operation(ao)
	return kind, nil
}
