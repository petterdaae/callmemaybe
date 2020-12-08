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
	KindInvalid   = memorymodel.ContextElementKindInvalid
	KindProcedure = memorymodel.ContextElementKindProcedure
	RAX           = assemblyoutput.RAX
	RBX           = assemblyoutput.RBX
	RDI           = assemblyoutput.RDI
	RSI           = assemblyoutput.RSI
	PRINTFORMAT64 = assemblyoutput.PRINTFORMAT64
)

func (exp ExpPlus) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (memorymodel.ContextElementKind, string, error) {
	kind, _, err := exp.Left.Generate(ao, mm)
	if err != nil {
		return KindInvalid, "", fmt.Errorf("failed to generate code for left side of plus exp: %w", err)
	}
	if kind != KindNumber {
		return KindInvalid, "", fmt.Errorf("can only add numbers")
	}

	mm.CurrentStackSize++
	ao.Push(RAX)

	kind, _, err = exp.Right.Generate(ao, mm)
	if err != nil {
		return KindInvalid, "", fmt.Errorf("failed to generate code for right side of plus exp: %w", err)
	}
	if kind != KindNumber {
		return KindInvalid, "", fmt.Errorf("can only add numbers")
	}

	mm.CurrentStackSize--
	ao.Pop(RBX)
	ao.Add(RAX, RBX)

	return KindNumber, "", nil
}

func (exp ExpMultiply) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (memorymodel.ContextElementKind, string, error) {
	kind, _, err := exp.Left.Generate(ao, mm)
	if err != nil {
		return KindInvalid, "", fmt.Errorf("failed to generate code for left side of multiply exp: %w", err)
	}
	if kind != KindNumber {
		return KindInvalid, "", fmt.Errorf("can only multiply numbers")
	}

	mm.CurrentStackSize++
	ao.Push(RAX)

	kind, _, err = exp.Right.Generate(ao, mm)
	if err != nil {
		return KindInvalid, "", fmt.Errorf("failed to generate code for right side of multiply exp: %w ", err)
	}
	if kind != KindNumber {
		return KindInvalid, "", fmt.Errorf("can only multiply numbers")
	}

	mm.CurrentStackSize--
	ao.Pop(RBX)
	ao.Imul(RAX, RBX)

	return KindNumber, "", nil
}

func (exp ExpParentheses) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (memorymodel.ContextElementKind, string, error) {
	return exp.Inside.Generate(ao, mm)
}

func (exp ExpNum) Generate(ao *assemblyoutput.AssemblyOutput, _ *memorymodel.MemoryModel) (memorymodel.ContextElementKind, string, error) {
	val := strconv.Itoa(exp.Value)
	ao.Mov(RAX, val)
	return KindNumber, "", nil
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

	if kind == KindNumber || kind == KindBool {
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

	if memorymodel.IsStackKind(kind) {
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

	if memorymodel.IsStackKind(kind) {
		return nil
	}

	return fmt.Errorf("only bools and ints are supported return types")
}

func (exp ExpFunction) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (memorymodel.ContextElementKind, string, error) {
	mm.PushNewContext(false)
	name := ao.PushProcedure(len(exp.Args), mm.CurrentStackSize, memorymodel.GetKindFromType(exp.ReturnType))

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
		if !memorymodel.IsStackKind(kind) {
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
	// TODO : implement
	return fmt.Errorf("not implemented")
}

func (expr ExpEquals) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (memorymodel.ContextElementKind, string, error) {
	kind, _, err := expr.Left.Generate(ao, mm)
	if err != nil {
		return KindInvalid, "", fmt.Errorf("failed to generate left expression of equals: %w", err)
	}
	if kind != KindNumber && kind != KindBool {
		return KindInvalid, "", fmt.Errorf("equals only supported for numns and bools")
	}

	mm.CurrentStackSize++
	ao.Push(RAX)

	kind, _, err = expr.Right.Generate(ao, mm)
	if err != nil {
		return KindInvalid, "", fmt.Errorf("failed to generate right expression of equals: %w", err)
	}
	if kind != KindNumber && kind != KindBool {
		return KindInvalid, "", fmt.Errorf("equals only supported for numns and bools")
	}

	mm.CurrentStackSize--
	ao.Pop(RBX)

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

	return KindBool, "", nil
}

func (expr ExpLess) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (memorymodel.ContextElementKind, string, error) {
	// TODO : implement
	return KindInvalid, "", fmt.Errorf("not implemented")
}

func (expr ExpGreater) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (memorymodel.ContextElementKind, string, error) {
	// TODO : implement
	return KindInvalid, "", fmt.Errorf("not implemented")
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
