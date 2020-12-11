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

func (exp ExpParentheses) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (memorymodel.ContextElementKind, string, error) {
	return exp.Inside.Generate(ao, mm)
}

func (exp ExpNum) Generate(ao *assemblyoutput.AssemblyOutput, _ *memorymodel.MemoryModel) (memorymodel.ContextElementKind, string, error) {
	val := strconv.Itoa(exp.Value)
	ao.Mov(RAX, val)
	return KindNumber, "", nil
}

func (exp ExpChar) Generate(ao *assemblyoutput.AssemblyOutput, _ *memorymodel.MemoryModel) (memorymodel.ContextElementKind, string, error) {
	rune := exp.Value[0]
	ao.Mov(RAX, fmt.Sprintf("%d", rune))
	return KindChar, "", nil
}

func (exp ExpIdentifier) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (memorymodel.ContextElementKind, string, error) {
	stackElement := mm.GetStackElement(exp.Name)
	if stackElement != nil {
		address := fmt.Sprintf("[rsp+%d]", (mm.CurrentStackSize-stackElement.StackSizeAfterPush)*8)
		ao.Mov(RAX, address)
		return stackElement.Kind, "", nil
	}

	procedureElement := mm.GetProcedureElement(exp.Name)
	if procedureElement != nil {
		return KindProcedure, procedureElement.Name, nil
	}

	return KindInvalid, "", fmt.Errorf("failed to find '%s' in current context", exp.Name)
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
	kind, name, err := stmt.Expression.Generate(ao, mm)
	if err != nil {
		return fmt.Errorf("failed to generate code for expression in assign statement: %w", err)
	}

	// Do not assign anything to placeholders
	if stmt.Identifier == "_" {
		return nil
	}

	if kind == KindNumber || kind == KindBool || kind == KindChar {
		mm.CurrentStackSize++
		ao.Push(RAX)
		mm.AddNameToCurrentStackElement(stmt.Identifier, kind)
	}

	if kind == KindProcedure {
		procedure := ao.GetProcedureByName(name)
		mm.AddProcedureAlias(name, stmt.Identifier, procedure.NumberOfArgs, procedure.ReturnKind)
	}

	return nil
}

func (stmt StmtPrintln) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) error {
	kind, _, err := stmt.Expression.Generate(ao, mm)
	if err != nil {
		return fmt.Errorf("failed to generate code for expression in println: %w", err)
	}

	if kind == KindChar {
		ao.Mov(RDI, PRINTCHARFORMAT)
		ao.Mov(RSI, RAX)
		ao.Xor(RAX, RAX)
		ao.CallPrintf()
		return nil
	}

	if memorymodel.IsIntOrBool(kind) {
		ao.Mov(RDI, PRINTFORMAT64)
		ao.Mov(RSI, RAX)
		ao.Xor(RAX, RAX)
		ao.CallPrintf()
		return nil
	}

	return fmt.Errorf("println only supports numbers and booleans")
}

func (stmt StmtReturn) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) error {
	kind, _, err := stmt.Expression.Generate(ao, mm)
	if err != nil {
		return fmt.Errorf("failed to evaluate expression when returning: %w", err)
	}

	procedure := ao.CurrentProcedure()
	for i := 0; i < mm.CurrentStackSize-procedure.StackSizeBeforeFunctionGeneration-1-procedure.NumberOfArgs; i++ {
		ao.Pop(RBX)
	}

	ao.Ret()

	if memorymodel.IsIntOrBool(kind) {
		return nil
	}

	return fmt.Errorf("only bools and ints are supported return types")
}

func (exp ExpFunction) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (memorymodel.ContextElementKind, string, error) {
	mm.PushNewContext(false)
	name := ao.PushProcedure(len(exp.Args), mm.CurrentStackSize, memorymodel.GetKindFromType(exp.ReturnType))

	mm.AddProcedureAlias(name, exp.Recurse, len(exp.Args), memorymodel.GetKindFromType(exp.ReturnType))

	initialStackSize := mm.CurrentStackSize

	for _, arg := range exp.Args {
		mm.CurrentStackSize++
		mm.AddNameToCurrentStackElement(arg.Identifier, memorymodel.GetKindFromType(arg.Type))
	}

	mm.CurrentStackSize++ // Return pointer is pushed to stack when calling procedure

	err := exp.Body.Generate(ao, mm)
	if err != nil {
		return KindInvalid, "", fmt.Errorf("failed to generate function body: %w", err)
	}

	mm.CurrentStackSize--

	for i := 0; i < mm.CurrentStackSize-initialStackSize-len(exp.Args); i++ {
		ao.Pop(RBX)
	}

	mm.CurrentStackSize = initialStackSize

	ao.Ret()

	mm.PopCurrentContext()
	ao.PopProcedure()

	return KindProcedure, name, nil
}

func (stmt FunctionCall) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (memorymodel.ContextElementKind, string, error) {
	for _, arg := range stmt.Arguments {
		kind, _, err := arg.Generate(ao, mm)
		if !memorymodel.IsIntOrBool(kind) {
			return KindInvalid, "", fmt.Errorf("only ints and bools are supported as argument types")
		}
		if err != nil {
			return KindInvalid, "", fmt.Errorf("failed to evaluate function argument: %w", err)
		}
		mm.CurrentStackSize++
		ao.Push(RAX)
	}

	procedureElement := mm.GetProcedureElement(stmt.Name)
	if procedureElement == nil {
		return KindInvalid, "", fmt.Errorf("no procedure with name '%s'", stmt.Name)
	}

	if procedureElement.NumberOfArgs != len(stmt.Arguments) {
		// TODO : also check types
		return KindInvalid, "", fmt.Errorf("mismatching number of arguments")
	}

	ao.Call(procedureElement.Name)

	for i := 0; i < len(stmt.Arguments); i++ {
		mm.CurrentStackSize--
		ao.Pop(RBX)
	}

	return procedureElement.ReturnKind, "", nil
}

func (stmt StmtIf) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) error {
	mm.PushNewContext(true)

	kind, _, err := stmt.Expression.Generate(ao, mm)
	if err != nil {
		return fmt.Errorf("failed to generate condition of id: %w", err)
	}
	if !memorymodel.IsIntOrBool(kind) {
		return fmt.Errorf("if conditions can only be stack kinds")
	}

	bodyStart := ao.GenerateUniqueName()
	bodyEnd := ao.GenerateUniqueName()

	ao.Cmp(RAX, "1")
	ao.Je(bodyStart)
	ao.Jne(bodyEnd)

	ao.NewSection(bodyStart)

	err = stmt.Body.Generate(ao, mm)
	if err != nil {
		return fmt.Errorf("failed to generate if body: %w", err)
	}

	ao.NewSection(bodyEnd)

	mm.PopCurrentContext()

	return nil
}

func (expr ExpBool) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (memorymodel.ContextElementKind, string, error) {
	var val string
	if expr.Value {
		val = "1"
	} else {
		val = "0"
	}
	ao.Mov(RAX, val)
	return KindBool, "", nil
}

func (expr ExpNegative) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (memorymodel.ContextElementKind, string, error) {
	kind, _, err := expr.Inside.Generate(ao, mm)
	if err != nil {
		return KindInvalid, "", fmt.Errorf("failed to generate expression inside negative: %w", err)
	}
	if !memorymodel.IsIntOrBool(kind) {
		return KindInvalid, "", fmt.Errorf("negative expressions only support stack kinds")
	}
	ao.Mov(RBX, RAX)
	ao.Mov(RAX, "0")
	ao.Sub(RAX, RBX)
	return KindNumber, "", nil
}
