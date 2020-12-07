package language

import (
	"fmt"
	"lang/language/assembly"
	"strconv"
)

func (exp ExpPlus) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	kind, err := exp.Left.Generate(gen)
	if err != nil {
		return InvalidExp, fmt.Errorf("failed to generate code for left side of plus exp: %w", err)
	}
	if kind != StackNumExp {
		return InvalidExp, fmt.Errorf("can only add numbers")
	}
	gen.push(rax)
	kind, err = exp.Right.Generate(gen)
	if err != nil {
		return InvalidExp, fmt.Errorf("failed to generate code for right side of plus exp: %w", err)
	}
	if kind != StackNumExp {
		return InvalidExp, fmt.Errorf("can only add numbers")
	}
	gen.pop(rbx)
	gen.add(rax, rbx)
	return StackNumExp, nil
}

func (exp ExpMultiply) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	kind, err := exp.Left.Generate(gen)
	if err != nil {
		return InvalidExp, fmt.Errorf("failed to generate code for left side of multiply exp: %w", err)
	}
	if kind != StackNumExp {
		return InvalidExp, fmt.Errorf("can only multiply numbers")
	}
	gen.push(rax)
	kind, err = exp.Right.Generate(gen)
	if err != nil {
		return InvalidExp, fmt.Errorf("failed to generate code for right side of multiply exp: %w ", err)
	}
	if kind != StackNumExp {
		return InvalidExp, fmt.Errorf("can only multiply numbers")
	}
	gen.pop(rbx)
	gen.mult(rax, rbx)
	return StackNumExp, nil
}

func (exp ExpParentheses) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	kind, err := exp.Inside.Generate(gen)
	if err != nil {
		return InvalidExp, fmt.Errorf("failed to generate code for inside of parentheses exp: %w", err)
	}
	return kind, nil
}

func (exp ExpNum) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	val := strconv.Itoa(exp.Value)
	gen.mov(rax, fmt.Sprintf("%s", val))
	return StackNumExp, nil
}

func (exp ExpIdentifier) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	kind, addr, err := gen.contexts.Get(exp.Name, gen.stackSize)
	if err != nil {
		return InvalidExp, fmt.Errorf("failed to evaluate identifer: %w", err)
	}

	if kind != assembly.ProcedureElem {
		gen.mov(rax, addr)
	}

	translatedKind := InvalidExp
	if kind == assembly.ProcedureElem {
		translatedKind = ProcExp
	} else if kind == assembly.StackBoolElem {
		translatedKind = StackBoolExp
	} else if kind == assembly.StackNumElem {
		translatedKind = StackNumExp
	}

	return translatedKind, nil
}

func (stmt StmtSeq) Generate(gen *AssemblyGenerator) error {
	for i := range stmt.Statements {
		err := stmt.Statements[i].Generate(gen)
		if err != nil {
			return fmt.Errorf("failed to generate code for statement in sequence: %w", err)
		}
	}
	return nil
}

func (stmt StmtAssign) Generate(gen *AssemblyGenerator) error {
	kind, err := stmt.Expression.Generate(gen)
	if err != nil {
		return fmt.Errorf("failed to generate code for expression in assign statement: %w", err)
	}
	if stmt.Identifier == "_" {
		return nil
	}
	if kind == StackNumExp || kind == StackBoolExp {
		contextKind := assembly.StackNumElem
		if kind == StackBoolExp {
			contextKind = assembly.StackBoolElem
		}
		operations, pushes, err := gen.contexts.StackInsert(stmt.Identifier, rax, gen.stackSize, contextKind)
		if err != nil {
			return fmt.Errorf("failed to insert element into stack: %w", err)
		}
		gen.AddOperations(operations)
		gen.stackSize += pushes
	}
	if kind == ProcExp {
		gen.NameNamelessProcedure(stmt.Identifier)
	}
	return nil
}

func (stmt StmtPrintln) Generate(gen *AssemblyGenerator) error {
	kind, err := stmt.Expression.Generate(gen)
	if err != nil {
		return fmt.Errorf("failed to generate code for expression in println: %w", err)
	}
	if kind != StackBoolExp && kind != StackNumExp {
		return fmt.Errorf("only num and bool printing is supported")
	}
	gen.println(rax)
	return nil
}

func (stmt StmtReturn) Generate(gen *AssemblyGenerator) error {
	kind, err := stmt.Expression.Generate(gen)
	if err != nil {
		return fmt.Errorf("failed to evaluate expression when returning: %w", err)
	}
	procedure := gen.contexts.GetTopProcedure()

	if procedure == nil {
		return fmt.Errorf("cant return outside a procedure")
	}

	if kind == EmptyExp {
		return fmt.Errorf("empty returns not allowed")
	}

	if !(procedure.ReturnType == "bool" && kind == StackBoolExp) && !(procedure.ReturnType == "int" && kind == StackNumExp) {
		return fmt.Errorf("mismatching types between function and return")
	}

	gen.jmp(fmt.Sprintf("%send", procedure.Name))
	return nil
}

func (exp ExpFunction) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	newContext, operations, pushes := gen.contexts.NewContext(true, gen.stackSize)
	gen.AddOperations(operations)
	gen.stackSize += pushes
	gen.contexts.Push(newContext)

	context := gen.contexts.Peek()
	context.Procedure.ReturnType = exp.ReturnType
	context.Procedure.NumberOfArgs = len(exp.Args)

	if len(exp.Args) != 0 {
		initStackSize := newContext.Procedure.StackSizeWhenInitialized
		for i := 1; i <= len(exp.Args); i++ {
			arg := exp.Args[i-1]
			gen.contexts.Peek().Stack[arg.Identifier] = initStackSize + i
		}
		gen.stackSize += len(exp.Args) + 1
	}

	err := exp.Body.Generate(gen)
	if err != nil {
		return InvalidExp, fmt.Errorf("failed to generate function body: %w", err)
	}

	if len(exp.Args) != 0 {
		gen.stackSize -= len(exp.Args) + 1
	}

	pops, context := gen.contexts.Pop(gen.stackSize)
	gen.stackSize -= pops
	gen.PushNamelessProcedure(context.Procedure)

	return ProcExp, nil
}

// TODO : require returns for typed functions
func (stmt FunctionCall) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	// Bind arguments
	for _, arg := range stmt.Arguments {
		arg.Generate(gen)
		gen.push(rax)
	}
	gen.stackSize -= len(stmt.Arguments)

	kind, actualName, err := gen.contexts.Get(stmt.Name, gen.stackSize)
	if kind != assembly.ProcedureElem || err != nil {
		return InvalidExp, fmt.Errorf("failed to call procedure with name: %s", stmt.Name)
	}

	var procedure *assembly.Procedure
	for i := range gen.AllProcedures {
		p := gen.AllProcedures[i]
		if p.Name == actualName {
			procedure = p
		}
	}
	if procedure == nil {
		return InvalidExp, fmt.Errorf("failed to find procedure with name: %s", actualName)
	}

	if procedure.NumberOfArgs != len(stmt.Arguments) {
		return InvalidExp, fmt.Errorf("mismatching number of arguments")
	}

	gen.mov(rcx, strconv.Itoa(gen.stackSize))
	gen.addOperation(fmt.Sprintf("call %s", actualName))

	// Pop bound arguments
	for range stmt.Arguments {
		gen.popWithoutDecreasingStackSize(rbx)
	}

	if err != nil {
		return InvalidExp, fmt.Errorf("failed to call function: %w", err)
	}

	if procedure.ReturnType == "empty" {
		return EmptyExp, nil
	}

	if procedure.ReturnType == "bool" {
		return StackBoolExp, nil
	}

	return StackNumExp, nil
}

func (stmt StmtIf) Generate(gen *AssemblyGenerator) error {
	// TODO
	return nil
}

func (expr ExpEquals) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	// TODO
	return InvalidExp, nil
}

func (expr ExpLess) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	// TODO
	return InvalidExp, nil
}

func (expr ExpGreater) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	// TODO
	return InvalidExp, nil
}

func (expr ExpBool) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	var val string
	if expr.Value {
		val = "1"
	} else {
		val = "0"
	}
	gen.mov(rax, fmt.Sprintf("%s", val))
	return StackBoolExp, nil
}
