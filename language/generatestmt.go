package language

import (
	"callmemaybe/language/assemblyoutput"
	"callmemaybe/language/memorymodel"
	"callmemaybe/language/typesystem"
	"fmt"
)

const (
	RAX                = assemblyoutput.RAX
	RBX                = assemblyoutput.RBX
	RDI                = assemblyoutput.RDI
	RSI                = assemblyoutput.RSI
	RDX                = assemblyoutput.RDX
	RCX                = assemblyoutput.RCX
)

func (stmt StmtSeq) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) error {
	for i := range stmt.Statements {
		err := stmt.Statements[i].Generate(ao, mm)
		if err != nil {
			return fmt.Errorf("statement in sequence: %w", err)
		}
	}
	return nil
}

func (stmt StmtAssign) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) error {
	kind, err := stmt.Expression.Generate(ao, mm)
	if err != nil {
		return fmt.Errorf("expression in assign: %w", err)
	}
	if stmt.Identifier == "_" {
		return nil
	}
	if kind.IsStorableOnStack() {
		if mm.Contains(stmt.Identifier) {
			member := mm.GetStackElement(stmt.Identifier)
			member.Type = kind
			ao.Mov(fmt.Sprintf("[rsp+%d]", (mm.CurrentStackSize-member.StackSizeAfterPush)*8), RAX)
		} else {
			mm.CurrentStackSize++
			ao.Push(RAX)
			mm.AddNameToCurrentStackElement(stmt.Identifier, kind)
		}
		return nil
	}
	return fmt.Errorf("expression in assign not storable on stack")
}

func (stmt StmtPrintln) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) error {
	kind, err := stmt.Expression.Generate(ao, mm)
	if err != nil {
		return fmt.Errorf("expression in println: %w", err)
	}
	if kind.RawType == typesystem.Char {
		ao.Mov(RBX, assemblyoutput.CharNewlineFormat)
		ao.Call(assemblyoutput.PrintRegisterWithFormat)
		return nil
	}
	if kind.RawType == typesystem.Int || kind.RawType == typesystem.Bool {
		ao.Mov(RBX, assemblyoutput.DigitNewlineFormat)
		ao.Call(assemblyoutput.PrintRegisterWithFormat)
		return nil
	}
	if kind.RawType == typesystem.List && kind.ListElementType.RawType == typesystem.Char {
		ao.Mov(RBX, assemblyoutput.CharFormat)
		ao.Call(assemblyoutput.PrintListWithFormat)
		ao.Mov(RAX, "10")
		ao.Mov(RBX, assemblyoutput.CharFormat)
		ao.Call(assemblyoutput.PrintRegisterWithFormat)
		return nil
	}
	return fmt.Errorf("unsupported type in println expression")
}

func (stmt StmtReturn) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) error {
	kind, err := stmt.Expression.Generate(ao, mm)
	if err != nil {
		return fmt.Errorf("return expression: %w", err)
	}
	procedure := ao.CurrentProcedure()
	if procedure == nil {
		return fmt.Errorf("returns are only allowed inside functions")
	}
	for i := 0; i < mm.CurrentStackSize-procedure.StackSizeBeforeFunctionGeneration-1-procedure.NumberOfArgs; i++ {
		ao.Pop(RBX)
	}
	ao.Ret()
	if !kind.IsPassable() {
		return fmt.Errorf("return type is not passable")
	}
	return nil
}

func (stmt StmtIf) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) error {
	mm.PushNewContext(true)
	initStacksize := mm.CurrentStackSize
	kind, err := stmt.Expression.Generate(ao, mm)
	if err != nil {
		return fmt.Errorf("if condition: %w", err)
	}
	if kind.RawType != typesystem.Bool {
		return fmt.Errorf("if condition is not a bool")
	}
	bodyStart := ao.GenerateUniqueName()
	bodyEnd := ao.GenerateUniqueName()
	ao.Cmp(RAX, "1")
	ao.Je(bodyStart)
	ao.Jne(bodyEnd)
	ao.NewSection(bodyStart)
	err = stmt.Body.Generate(ao, mm)
	if err != nil {
		return fmt.Errorf("if body: %w", err)
	}
	for i := 0; i < mm.CurrentStackSize-initStacksize; i++ {
		mm.CurrentStackSize--
		ao.Pop(RBX)
	}
	ao.NewSection(bodyEnd)
	mm.PopCurrentContext()
	return nil
}

func (stmt StmtLoop) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) error {
	mm.PushNewContext(true)
	initStackSize := mm.CurrentStackSize
	bodyStart := ao.GenerateUniqueName()
	conditionStart := ao.GenerateUniqueName()
	ao.Jmp(conditionStart)
	ao.NewSection(bodyStart)
	err := stmt.Body.Generate(ao, mm)
	if err != nil {
		return fmt.Errorf("loop body: %w", err)
	}
	for i := 0; i < mm.CurrentStackSize-initStackSize; i++ {
		ao.Pop(RBX)
	}
	mm.CurrentStackSize = initStackSize
	ao.NewSection(conditionStart)
	kind, err := stmt.Condition.Generate(ao, mm)
	if err != nil {
		return fmt.Errorf("loop condition: %w", err)
	}
	if kind.RawType != typesystem.Bool {
		return fmt.Errorf("loop condition is not bool")
	}
	ao.Cmp(RAX, "1")
	ao.Je(bodyStart)
	mm.PopCurrentContext()
	return nil
}

func (stmt StmtStructDeclaration) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) error {
	mm.NewStructType(stmt.Type.StructName, stmt.Type)
	return nil
}

func (stmt StmtUpdateList) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) error {
	listKind, err := stmt.List.Generate(ao, mm)
	if err != nil {
		return fmt.Errorf("failed to generate list in update stmt")
	}
	if listKind.RawType != typesystem.List {
		return fmt.Errorf("expected list kind")
	}
	mm.CurrentStackSize++
	ao.Push(RAX)

	newValueKind, err := stmt.NewValue.Generate(ao, mm)
	if err != nil {
		return fmt.Errorf("failed to generate new value in update stmt")
	}
	if !newValueKind.Equals(*listKind.ListElementType) {
		return fmt.Errorf("new value type does not match list element type")
	}
	mm.CurrentStackSize++
	ao.Push(RAX)

	indexKind, err := stmt.Index.Generate(ao, mm)
	if err != nil {
		return fmt.Errorf("faied to generate index in update stmt")
	}
	if indexKind.RawType != typesystem.Int {
		return fmt.Errorf("index in update stmt was not int")
	}
	mm.CurrentStackSize++
	ao.Push(RAX)

	mm.CurrentStackSize -= 3
	ao.Pop(RAX)
	ao.Pop(RBX)
	ao.Pop(RCX)

	ao.Mov(fmt.Sprintf("qword [%s*8+%s+8]", RAX, RCX), RBX)
	return nil
}

func (stmt StmtUpdateStruct) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) error {
	structKind, err := stmt.Struct.Generate(ao, mm)
	if err != nil {
		return fmt.Errorf("failed to generate struct in update stmt")
	}
	if structKind.RawType != typesystem.Struct {
		return fmt.Errorf("expected struct kind")
	}
	mm.CurrentStackSize++
	ao.Push(RAX)

	newValueKind, err := stmt.NewValue.Generate(ao, mm)
	if err != nil {
		return fmt.Errorf("failed to generate new value in update stmt")
	}
	mm.CurrentStackSize++
	ao.Push(RAX)

	i := 0
	found := false
	for _, field := range structKind.StructMembers {
		if field.Name == stmt.Member {
			if !newValueKind.Equals(field.Type) {
				return fmt.Errorf("wrong type in update struct stmt")
			}
			found = true
			break
		}
		i++
	}

	if !found {
		return fmt.Errorf("%s is not a member of this struct", stmt.Member)
	}

	mm.CurrentStackSize -= 2
	ao.Pop(RAX)
	ao.Pop(RBX)

	ao.Mov(fmt.Sprintf("[%s+%d]", RBX, i*8), RAX)

	return nil
}

