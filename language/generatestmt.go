package language

import (
	"callmemaybe/language/assemblyoutput"
	"callmemaybe/language/common"
	"callmemaybe/language/memorymodel"
	"fmt"
)

const (
	KindNumber    = common.ContextElementKindNumber
	KindBool      = common.ContextElementKindBoolean
	KindChar      = common.ContextElementKindChar
	KindInvalid   = common.ContextElementKindInvalid
	KindList      = common.ContextElementKindListReference
	KindProcedure = common.ContextElementKindProcedure
	RAX           = assemblyoutput.RAX
	RBX           = assemblyoutput.RBX
	RDI           = assemblyoutput.RDI
	RSI           = assemblyoutput.RSI
	RDX           = assemblyoutput.RDX
	RCX           = assemblyoutput.RCX
	PRINTFORMAT64 = assemblyoutput.PRINTFORMAT64
	PRINTCHARFORMAT = assemblyoutput.PRINTCHARFORMAT
)





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
		mm.AddNameToCurrentStackElement(stmt.Identifier, result.Kind, result.ListElementKind)
	}

	if result.Kind == KindProcedure {
		procedure := ao.GetProcedureByName(result.ProcedureName)
		mm.AddProcedureAlias(result.ProcedureName, stmt.Identifier, procedure.NumberOfArgs, procedure.ReturnKind, procedure.Args)
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

	if result.Kind == KindNumber || result.Kind == KindBool {
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

	if result.Kind.IsStoredOnStack() {
		return nil
	}

	return fmt.Errorf("only bools and ints are supported return types")
}

func (stmt StmtIf) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) error {
	mm.PushNewContext(true)

	result := stmt.Expression.Generate(ao, mm)
	if result.Error != nil {
		return fmt.Errorf("failed to generate condition of id: %w", result.Error)
	}
	if result.Kind != KindBool {
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

