package language

import (
	"callmemaybe/language/assemblyoutput"
	"callmemaybe/language/memorymodel"
	"callmemaybe/language/typesystem"
	"fmt"
	"strconv"
)

func (exp ExpParentheses) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (typesystem.Type, error) {
	return exp.Inside.Generate(ao, mm)
}

func (exp ExpNum) Generate(ao *assemblyoutput.AssemblyOutput, _ *memorymodel.MemoryModel) (typesystem.Type, error) {
	val := strconv.Itoa(exp.Value)
	ao.Mov(RAX, val)
	return typesystem.NewInt(), nil
}

func (exp ExpChar) Generate(ao *assemblyoutput.AssemblyOutput, _ *memorymodel.MemoryModel) (typesystem.Type, error) {
	rune := exp.Value[0]
	ao.Mov(RAX, fmt.Sprintf("%d", rune))
	return typesystem.NewChar(), nil
}

func (exp ExpIdentifier) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (typesystem.Type, error) {
	stackElement := mm.GetStackElement(exp.Name)
	if stackElement != nil {
		address := fmt.Sprintf("[rsp+%d]", (mm.CurrentStackSize-stackElement.StackSizeAfterPush)*8)
		ao.Mov(RAX, address)
		return stackElement.Type, nil
	}
	return typesystem.NewInvalid(), fmt.Errorf("missing from context: %s", exp.Name)
}

func (exp ExpFunction) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (typesystem.Type, error) {
	mm.PushNewContext(false)
	name := ao.PushProcedure(mm.CurrentStackSize, len(exp.Type.FunctionArgumentTypes))
	initialStackSize := mm.CurrentStackSize

	argNames := make(map[string]bool)
	for _, arg := range exp.Type.FunctionArgumentTypes {
		_, exists := argNames[arg.Name]
		if exists {
			return typesystem.NewInvalid(), fmt.Errorf("argument names should be unique")
		}
		if arg.Type.RawType == typesystem.Struct {
			var ok bool
			arg.Type, ok = mm.GetStructType(arg.Type.StructName)
			if !ok {
				return typesystem.NewInvalid(), fmt.Errorf("type does not exist")
			}
		}
		mm.CurrentStackSize++
		mm.AddNameToCurrentStackElement(arg.Name, arg.Type)
		argNames[arg.Name] = true
	}
	if len(argNames) != len(exp.Type.FunctionArgumentTypes) {
		return typesystem.NewInvalid(), fmt.Errorf("mismatching number of arguments")
	}

	mm.CurrentStackSize++

	ao.Mov(RAX, name)
	mm.CurrentStackSize++
	ao.Push(RAX)
	mm.AddNameToCurrentStackElement(exp.Recurse, exp.Type)

	err := exp.Body.Generate(ao, mm)
	if err != nil {
		return typesystem.NewInvalid(), fmt.Errorf("function body: %w", err)
	}

	mm.CurrentStackSize--

	for i := 0; i < mm.CurrentStackSize-initialStackSize-len(exp.Type.FunctionArgumentTypes); i++ {
		ao.Pop(RBX)
	}
	mm.CurrentStackSize = initialStackSize
	ao.Ret()
	mm.PopCurrentContext()
	ao.PopProcedure()
	ao.Mov(RAX, name)
	return exp.Type, nil
}

func (stmt FunctionCall) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (typesystem.Type, error) {
	kind, err := stmt.Exp.Generate(ao, mm)

	mm.CurrentStackSize++
	ao.Push(RAX)

	if err != nil {
		return typesystem.NewInvalid(), fmt.Errorf("call expression: %w", err)
	}
	if kind.RawType != typesystem.Function {
		return typesystem.NewInvalid(), fmt.Errorf("can only call functions")
	}

	if len(kind.FunctionArgumentTypes) != len(stmt.Arguments) {
		return typesystem.NewInvalid(), fmt.Errorf("mismathcing number of arguments in call")
	}

	for i := 0; i < len(stmt.Arguments); i++ {
		_kind, err := stmt.Arguments[i].Generate(ao, mm)
		if err != nil {
			return typesystem.NewInvalid(), fmt.Errorf("argument in call: %w", err)
		}
		if !_kind.IsPassable() {
			return typesystem.NewInvalid(), fmt.Errorf("argument type must be passable")
		}

		mm.CurrentStackSize++
		ao.Push(RAX)
		argKind := kind.FunctionArgumentTypes[i].Type
		if !argKind.Equals(_kind) {
			return typesystem.NewInvalid(), fmt.Errorf("mismatching argument types in call")
		}
	}

	ao.Call(fmt.Sprintf("[rsp+%d]", len(stmt.Arguments)*8))
	mm.CurrentStackSize--
	ao.Pop(RBX)

	for i := 0; i < len(stmt.Arguments); i++ {
		mm.CurrentStackSize--
		ao.Pop(RBX)
	}

	if kind.FunctionReturnType == nil {
		return typesystem.NewInvalid(), fmt.Errorf("functionreturntype is nil")
	}

	if kind.FunctionReturnType.RawType == typesystem.Struct {
		val, ok := mm.GetStructType(kind.FunctionReturnType.StructName)
		if !ok {
			return typesystem.NewInvalid(), fmt.Errorf("invalid struct type")
		}
		return val, nil
	}

	return *kind.FunctionReturnType, nil
}

func (expr ExpBool) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (typesystem.Type, error) {
	var val string
	if expr.Value {
		val = "1"
	} else {
		val = "0"
	}
	ao.Mov(RAX, val)
	return typesystem.NewBool(), nil
}

func (expr ExpNegative) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (typesystem.Type, error) {
	kind, err := expr.Inside.Generate(ao, mm)
	if err != nil {
		return typesystem.NewInvalid(), fmt.Errorf("negative: %w", err)
	}
	if !kind.IsAlgebraic() {
		return typesystem.NewInvalid(), fmt.Errorf("negative expressions only support algebraic kinds")
	}
	ao.Mov(RBX, RAX)
	ao.Mov(RAX, "0")
	ao.Sub(RAX, RBX)
	return typesystem.NewInt(), nil
}

func (expr ExpList) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (typesystem.Type, error) {
	if expr.Size < 1 {
		return typesystem.NewInvalid(), fmt.Errorf("the size of a list must be a positive number")
	}
	ao.Mov(RDI, fmt.Sprintf("%d", 8*(expr.Size+1)))
	ao.Call("malloc")
	ao.Mov(RDX, RAX)

	if expr.Size < len(expr.Elements) {
		return typesystem.NewInvalid(), fmt.Errorf("too many elements in list")
	}

	ao.Mov(fmt.Sprintf("qword [%s]", RDX), fmt.Sprintf("%d", expr.Size))
	for i, element := range expr.Elements {
		mm.CurrentStackSize++
		ao.Push(RDX)
		kind, err := element.Generate(ao, mm)
		if err != nil {
			return typesystem.NewInvalid(), fmt.Errorf("failed to generate expression of element in list: %w", err)
		}
		if !kind.Equals(*expr.Type.ListElementType) {
			return typesystem.NewInvalid(), fmt.Errorf("list element has the wrong type")
		}
		mm.CurrentStackSize--
		ao.Pop(RDX)
		ao.Mov(fmt.Sprintf("qword [%s+%d]", RDX, (i+1)*8), RAX)
	}

	ao.Mov(RAX, RDX)

	return expr.Type, nil
}

func (expr ExpGetFromList) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (typesystem.Type, error) {
	kind, err := expr.Index.Generate(ao, mm)
	if err != nil {
		return typesystem.NewInvalid(), fmt.Errorf("failed to evaluate index")
	}
	if kind.RawType != typesystem.Int {
		return typesystem.NewInvalid(), fmt.Errorf("only integers are valid indexes")
	}
	mm.CurrentStackSize++
	ao.Push(RAX)
	kind, err = expr.List.Generate(ao, mm)
	if err != nil {
		return typesystem.NewInvalid(), fmt.Errorf("failed to generate code for list in get expression: %w", err)
	}
	if kind.RawType != typesystem.List {
		return typesystem.NewInvalid(), fmt.Errorf("can only get from lists by index")
	}
	mm.CurrentStackSize--
	ao.Pop(RCX)
	ao.Mov(RDX, RAX)
	ao.Mov(RAX, fmt.Sprintf("[rdx+8*%s+8]", RCX))

	if kind.ListElementType == nil {
		return typesystem.Type{}, fmt.Errorf("listelementtype is nil")
	}

	return *kind.ListElementType, nil
}

func (expr StructExp) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (typesystem.Type, error) {
	typeFromMemoryModel, ok := mm.GetStructType(expr.Name)
	if !ok {
		return typesystem.NewInvalid(), fmt.Errorf("struct type does not exist")
	}
	_type := typesystem.Type{
		RawType:               typesystem.Struct,
		StructName:            expr.Name,
	}
	ao.Mov(RDI, fmt.Sprintf("%d", 8*len(expr.Members)))
	ao.Call("malloc")
	ao.Mov(RDX, RAX)
	i := 0
	for _, member := range expr.Members {
		actual := typeFromMemoryModel.StructMembers[i]
		if actual.Name != member.Name {
			return typesystem.NewInvalid(), fmt.Errorf("invalid struct field name")
		}
		mm.CurrentStackSize++
		ao.Push(RDX)
		kind, err := member.Exp.Generate(ao, mm)
		mm.CurrentStackSize--
		ao.Pop(RDX)
		if !kind.Equals(actual.Type) {
			return typesystem.NewInvalid(), fmt.Errorf("invalid field type")
		}
		if err != nil {
			return typesystem.NewInvalid(), fmt.Errorf("expression in struct init: %w", err)
		}
		ao.Mov(fmt.Sprintf("qword [%s+%d]", RDX, i*8), RAX)
		_type.StructMembers = append(_type.StructMembers, typesystem.NamedType{
			Name: member.Name,
			Type: kind,
		})
		i += 1
	}

	if len(expr.Members) != len(typeFromMemoryModel.StructMembers) {
		return typesystem.NewInvalid(), fmt.Errorf("mismatching number of arguments in field declaration")
	}

	mm.CurrentStackSize++
	ao.Mov(RAX, RDX)
	ao.Push(RAX)

	return _type, nil
}

func (expr ExpReadFromStruct) Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) (typesystem.Type, error) {
	kind, err := expr.Struct.Generate(ao, mm)
	ao.Mov(RDX, RAX)
	if err != nil {
		return typesystem.NewInvalid(), fmt.Errorf("struct in read from struct: %w", err)
	}
	if kind.RawType != typesystem.Struct {
		return typesystem.NewInvalid(), fmt.Errorf("can only read from structs")
	}

	i := 0
	for _, member := range kind.StructMembers {
		if member.Name == expr.Field {
			ao.Mov(RAX, fmt.Sprintf("[%s+%d]", RDX, i*8))
			return member.Type, nil
		}
		i++
	}

	return typesystem.NewInvalid(), fmt.Errorf("invalid field")
}
