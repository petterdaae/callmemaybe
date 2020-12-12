package language

import (
	"callmemaybe/language/assemblyoutput"
	"callmemaybe/language/memorymodel"
	"fmt"
	"strconv"
)

func (exp ExpParentheses) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerateResult {
	return exp.Inside.Generate(ao, mm)
}

func (exp ExpNum) Generate(ao *assemblyoutput.AssemblyOutput, _ *memorymodel.MemoryModel) GenerateResult {
	val := strconv.Itoa(exp.Value)
	ao.Mov(RAX, val)
	return NumberResult()
}

func (exp ExpChar) Generate(ao *assemblyoutput.AssemblyOutput, _ *memorymodel.MemoryModel) GenerateResult {
	rune := exp.Value[0]
	ao.Mov(RAX, fmt.Sprintf("%d", rune))
	return CharResult()
}

func (exp ExpIdentifier) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerateResult {
	stackElement := mm.GetStackElement(exp.Name)
	if stackElement != nil {
		address := fmt.Sprintf("[rsp+%d]", (mm.CurrentStackSize-stackElement.StackSizeAfterPush)*8)
		ao.Mov(RAX, address)
		return CustomResult(stackElement.Kind, stackElement.ListElementKind, "", nil)
	}

	procedureElement := mm.GetProcedureElement(exp.Name)
	if procedureElement != nil {
		return ProcedureResult(procedureElement.Name)
	}

	return ErrorResult(fmt.Errorf("failed to find '%s' in current context", exp.Name))
}

func (exp ExpFunction) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerateResult {
	mm.PushNewContext(false)
	name := ao.PushProcedure(len(exp.Args), mm.CurrentStackSize, exp.ReturnType, exp.Args)

	if exp.Recurse != "" {
		mm.AddProcedureAlias(name, exp.Recurse, len(exp.Args), exp.ReturnType, exp.Args)
	}

	initialStackSize := mm.CurrentStackSize
	argNames := make(map[string]bool)

	for _, arg := range exp.Args {
		mm.CurrentStackSize++
		mm.AddNameToCurrentStackElement(arg.Identifier, arg.Type, KindInvalid)
		argNames[arg.Identifier] = true
	}

	if len(argNames) != len(exp.Args) {
		return ErrorResult(fmt.Errorf("argument names should be unique"))
	}

	mm.CurrentStackSize++

	err := exp.Body.Generate(ao, mm)
	if err != nil {
		return ErrorResult(fmt.Errorf("failed to generate function body: %w", err))
	}

	mm.CurrentStackSize--

	for i := 0; i < mm.CurrentStackSize-initialStackSize-len(exp.Args); i++ {
		ao.Pop(RBX)
	}

	mm.CurrentStackSize = initialStackSize
	ao.Ret()
	mm.PopCurrentContext()
	ao.PopProcedure()

	return ProcedureResult(name)
}

func (stmt FunctionCall) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerateResult {
	procedureElement := mm.GetProcedureElement(stmt.Name)
	if procedureElement == nil {
		return ErrorResult(fmt.Errorf("no procedure with name '%s'", stmt.Name))
	}

	if procedureElement.FunctionNumberOfArgs != len(stmt.Arguments) {
		return ErrorResult(fmt.Errorf("mismatching number of arguments"))
	}

	for i := 0; i < procedureElement.FunctionNumberOfArgs; i++ {
		result := stmt.Arguments[i].Generate(ao, mm)
		if !result.Kind.IsPassable() {
			return ErrorResult(fmt.Errorf("argument kind is not passable"))
		}
		if result.IsError() {
			return result.WrapError("failed to evaluate function argument")
		}
		mm.CurrentStackSize++
		ao.Push(RAX)
		argKind := procedureElement.FunctionArguments[i].Type
		if argKind != result.Kind {
			return ErrorResult(fmt.Errorf("mismatching argument types when calling function"))
		}
	}

	ao.Call(procedureElement.Name)

	for i := 0; i < len(stmt.Arguments); i++ {
		mm.CurrentStackSize--
		ao.Pop(RBX)
	}

	return CustomResult(procedureElement.FunctionReturnKind, KindInvalid, "", nil)
}

func (expr ExpBool) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerateResult {
	var val string
	if expr.Value {
		val = "1"
	} else {
		val = "0"
	}
	ao.Mov(RAX, val)
	return BoolResult()
}

func (expr ExpNegative) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerateResult {
	result := expr.Inside.Generate(ao, mm)
	if result.Error != nil {
		return ErrorResult(fmt.Errorf("failed to generate expression inside negative: %w", result.Error))
	}
	if !result.Kind.IsAlgebraic() {
		return ErrorResult(fmt.Errorf("negative expressions only support algebraic kinds"))
	}
	ao.Mov(RBX, RAX)
	ao.Mov(RAX, "0")
	ao.Sub(RAX, RBX)
	return NumberResult()
}

func (expr ExpList) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerateResult {
	if !expr.Type.IsStoredOnStack() {
		return ErrorResult(fmt.Errorf("list currently only support ints, bools and chars"))
	}
	if expr.Size < 1 {
		return ErrorResult(fmt.Errorf("the size of a list must be a positive number"))
	}
	ao.Mov(RDI, fmt.Sprintf("%d", 8*expr.Size))
	ao.Call("malloc")
	ao.Mov(RDX, RAX)

	if expr.Size > len(expr.Elements) {
		return ErrorResult(fmt.Errorf("too many elements in list"))
	}

	for i, element := range expr.Elements {
		result := element.Generate(ao, mm)
		if result.Error != nil {
			return ErrorResult(fmt.Errorf("failed to generate expression of element in list: %w", result.Error))
		}
		if result.Kind != expr.Type {
			return ErrorResult(fmt.Errorf("list element has the wrong type"))
		}
		ao.Mov(fmt.Sprintf("qword [%s+%d]", RDX, i*8), RAX)
	}

	ao.Mov(RAX, RDX)

	return ListResult(expr.Type)
}

func (expr ExpGetFromList) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerateResult {
	result := expr.Index.Generate(ao, mm)
	if result.Error != nil {
		return ErrorResult(fmt.Errorf("failed to evaluate index"))
	}
	if result.Kind != KindNumber {
		return ErrorResult(fmt.Errorf("only integers are valid indexes"))
	}
	mm.CurrentStackSize++
	ao.Push(RAX)
	result = expr.List.Generate(ao, mm)
	if result.Error != nil {
		return ErrorResult(fmt.Errorf("failed to generate code for list in get expression: %w", result.Error))
	}
	if result.Kind != KindList {
		return ErrorResult(fmt.Errorf("can only get from lists by index"))
	}
	mm.CurrentStackSize--
	ao.Pop(RCX)
	ao.Mov(RDX, RAX)
	ao.Mov(RAX, fmt.Sprintf("[rdx+8*%s]", RCX))

	return CustomResult(result.ListElementKind, KindInvalid, "", nil)
}
