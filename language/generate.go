package language

import (
	"fmt"
	"lang/language/assemblyoutput"
	"lang/language/memorymodel"
	"strconv"
)

const (
	KindNumber    = memorymodel.ContextElementKindNumber
	KindBool      = memorymodel.ContextElementKindBoolean
	KindChar      = memorymodel.ContextElementKindChar
	KindInvalid   = memorymodel.ContextElementKindInvalid
	KindList      = memorymodel.ContextElementKindListReference
	KindProcedure = memorymodel.ContextElementKindProcedure
	RAX           = assemblyoutput.RAX
	RBX           = assemblyoutput.RBX
	RDI           = assemblyoutput.RDI
	RSI           = assemblyoutput.RSI
	RDX           = assemblyoutput.RDX
	RCX           = assemblyoutput.RCX
	PRINTFORMAT64 = assemblyoutput.PRINTFORMAT64
	PRINTCHARFORMAT = assemblyoutput.PRINTCHARFORMAT
)

func (exp ExpParentheses) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerationResult {
	return exp.Inside.Generate(ao, mm)
}

func (exp ExpNum) Generate(ao *assemblyoutput.AssemblyOutput, _ *memorymodel.MemoryModel) GenerationResult {
	val := strconv.Itoa(exp.Value)
	ao.Mov(RAX, val)
	return NumberKind()
}

func (exp ExpChar) Generate(ao *assemblyoutput.AssemblyOutput, _ *memorymodel.MemoryModel) GenerationResult {
	rune := exp.Value[0]
	ao.Mov(RAX, fmt.Sprintf("%d", rune))
	return CharKind()
}

func (exp ExpIdentifier) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerationResult {
	stackElement := mm.GetStackElement(exp.Name)
	if stackElement != nil {
		address := fmt.Sprintf("[rsp+%d]", (mm.CurrentStackSize-stackElement.StackSizeAfterPush)*8)
		ao.Mov(RAX, address)
		return CustomKind(stackElement.Kind, "", nil)
	}

	procedureElement := mm.GetProcedureElement(exp.Name)
	if procedureElement != nil {
		return ProcedureKind(procedureElement.Name)
	}

	return ErrorKind(fmt.Errorf("failed to find '%s' in current context", exp.Name))
}

func (stmt StmtSeq) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) error {
	for i := range stmt.Statements {
		err := stmt.Statements[i].Generate(ao, mm)
		if err != nil {
			return fmt.Errorf("failed to generate code for statement in sequence: %w", err)
		}
	}
	return nil
}

func (stmt StmtAssign) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) error {
	result := stmt.Expression.Generate(ao, mm)
	if result.Error != nil {
		return fmt.Errorf("failed to generate code for expression in assign statement: %w", result.Error)
	}

	// Do not assign anything to placeholders
	if stmt.Identifier == "_" {
		return nil
	}

	if result.Kind == KindNumber || result.Kind == KindBool || result.Kind == KindChar || result.Kind == KindList {
		mm.CurrentStackSize++
		ao.Push(RAX)
		mm.AddNameToCurrentStackElement(stmt.Identifier, result.Kind)
	}

	if result.Kind == KindProcedure {
		procedure := ao.GetProcedureByName(result.ProcedureName)
		mm.AddProcedureAlias(result.ProcedureName, stmt.Identifier, procedure.NumberOfArgs, procedure.ReturnKind)
	}

	return nil
}

func (stmt StmtPrintln) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) error {
	result := stmt.Expression.Generate(ao, mm)
	if result.Error != nil {
		return fmt.Errorf("failed to generate code for expression in println: %w", result.Error)
	}

	if result.Kind == KindChar {
		ao.Mov(RDI, PRINTCHARFORMAT)
		ao.Mov(RSI, RAX)
		ao.Xor(RAX, RAX)
		ao.CallPrintf()
		return nil
	}

	if memorymodel.IsIntOrBool(result.Kind) {
		ao.Mov(RDI, PRINTFORMAT64)
		ao.Mov(RSI, RAX)
		ao.Xor(RAX, RAX)
		ao.CallPrintf()
		return nil
	}

	return fmt.Errorf("println only supports numbers and booleans")
}

func (stmt StmtReturn) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) error {
	result := stmt.Expression.Generate(ao, mm)
	if result.Error != nil {
		return fmt.Errorf("failed to evaluate expression when returning: %w", result.Error)
	}

	procedure := ao.CurrentProcedure()
	for i := 0; i < mm.CurrentStackSize-procedure.StackSizeBeforeFunctionGeneration-1-procedure.NumberOfArgs; i++ {
		ao.Pop(RBX)
	}

	ao.Ret()

	if memorymodel.IsIntOrBool(result.Kind) {
		return nil
	}

	return fmt.Errorf("only bools and ints are supported return types")
}

func (exp ExpFunction) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerationResult {
	mm.PushNewContext(false)
	name := ao.PushProcedure(len(exp.Args), mm.CurrentStackSize, exp.ReturnType)

	mm.AddProcedureAlias(name, exp.Recurse, len(exp.Args), exp.ReturnType)

	initialStackSize := mm.CurrentStackSize

	for _, arg := range exp.Args {
		mm.CurrentStackSize++
		mm.AddNameToCurrentStackElement(arg.Identifier, arg.Type)
	}

	mm.CurrentStackSize++ // Return pointer is pushed to stack when calling procedure

	err := exp.Body.Generate(ao, mm)
	if err != nil {
		return ErrorKind(fmt.Errorf("failed to generate function body: %w", err))
	}

	mm.CurrentStackSize--

	for i := 0; i < mm.CurrentStackSize-initialStackSize-len(exp.Args); i++ {
		ao.Pop(RBX)
	}

	mm.CurrentStackSize = initialStackSize

	ao.Ret()

	mm.PopCurrentContext()
	ao.PopProcedure()

	return ProcedureKind(name)
}

func (stmt FunctionCall) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerationResult {
	for _, arg := range stmt.Arguments {
		result := arg.Generate(ao, mm)
		if !memorymodel.IsIntOrBool(result.Kind) {
			return ErrorKind(fmt.Errorf("only ints and bools are supported as argument types"))
		}
		if result.Error != nil {
			return ErrorKind(fmt.Errorf("failed to evaluate function argument: %w", result.Error))
		}
		mm.CurrentStackSize++
		ao.Push(RAX)
	}

	procedureElement := mm.GetProcedureElement(stmt.Name)
	if procedureElement == nil {
		return ErrorKind(fmt.Errorf("no procedure with name '%s'", stmt.Name))
	}

	if procedureElement.NumberOfArgs != len(stmt.Arguments) {
		// TODO : also check types
		return ErrorKind(fmt.Errorf("mismatching number of arguments"))
	}

	ao.Call(procedureElement.Name)

	for i := 0; i < len(stmt.Arguments); i++ {
		mm.CurrentStackSize--
		ao.Pop(RBX)
	}

	return CustomKind(procedureElement.ReturnKind, "", nil)
}

func (stmt StmtIf) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) error {
	mm.PushNewContext(true)

	result := stmt.Expression.Generate(ao, mm)
	if result.Error != nil {
		return fmt.Errorf("failed to generate condition of id: %w", result.Error)
	}
	if !memorymodel.IsIntOrBool(result.Kind) {
		return fmt.Errorf("if conditions can only be stack kinds")
	}

	bodyStart := ao.GenerateUniqueName()
	bodyEnd := ao.GenerateUniqueName()

	ao.Cmp(RAX, "1")
	ao.Je(bodyStart)
	ao.Jne(bodyEnd)

	ao.NewSection(bodyStart)

	err := stmt.Body.Generate(ao, mm)
	if err != nil {
		return fmt.Errorf("failed to generate if body: %w", err)
	}

	ao.NewSection(bodyEnd)

	mm.PopCurrentContext()

	return nil
}

func (expr ExpBool) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerationResult {
	var val string
	if expr.Value {
		val = "1"
	} else {
		val = "0"
	}
	ao.Mov(RAX, val)
	return BoolKind()
}

func (expr ExpNegative) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerationResult {
	result := expr.Inside.Generate(ao, mm)
	if result.Error != nil {
		return ErrorKind(fmt.Errorf("failed to generate expression inside negative: %w", result.Error))
	}
	if !memorymodel.IsIntOrBool(result.Kind) {
		return ErrorKind(fmt.Errorf("negative expressions only support stack kinds"))
	}
	ao.Mov(RBX, RAX)
	ao.Mov(RAX, "0")
	ao.Sub(RAX, RBX)
	return NumberKind()
}

func (expr ExpList) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerationResult {
	if !memorymodel.IsIntOrBool(expr.Type) && expr.Type != memorymodel.ContextElementKindChar {
		return ErrorKind(fmt.Errorf("list currently only support ints, bools and chars"))
	}
	if expr.Size < 1 {
		return ErrorKind(fmt.Errorf("the size of a list must be a positive number"))
	}
	ao.Mov(RDI, fmt.Sprintf("%d", 8 * expr.Size))
	ao.Call("malloc")
	ao.Mov(RDX, RAX)

	if expr.Size > len(expr.Elements) {
		return ErrorKind(fmt.Errorf("too many elements in list"))
	}

	for i, element := range expr.Elements {
		result := element.Generate(ao, mm)
		if result.Error != nil {
			return ErrorKind(fmt.Errorf("failed to generate expression of element in list: %w", result.Error))
		}
		if result.Kind != expr.Type {
			return ErrorKind(fmt.Errorf("list element has the wrong type"))
		}
		ao.Mov(fmt.Sprintf("qword [%s+%d]", RDX, i * 8), RAX)
	}

	ao.Mov(RAX, RDX)

	return ListKind(expr.Type)
}

func (expr ExpGetFromList) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerationResult {
	result := expr.Index.Generate(ao, mm)
	if result.Error != nil {
		return ErrorKind(fmt.Errorf("failed to evaluate index"))
	}
	if result.Kind != KindNumber {
		return ErrorKind(fmt.Errorf("only integers are valid indexes"))
	}
	mm.CurrentStackSize++
	ao.Push(RAX)
	result = expr.List.Generate(ao, mm)
	if result.Error != nil {
		return ErrorKind(fmt.Errorf("failed to generate code for list in get expression: %w", result.Error))
	}
	if result.Kind != KindList {
		return ErrorKind(fmt.Errorf("can only get from lists by index"))
	}
	mm.CurrentStackSize--
	ao.Pop(RCX)
	ao.Mov(RDX, RAX)
	ao.Mov(RAX, fmt.Sprintf("[rdx+8*%s]", RCX))

	return CustomKind(result.ListElementKind, "", nil)
}
