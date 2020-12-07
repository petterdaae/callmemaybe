package language

import (
	"fmt"
	"lang/language/assembly"
	"strconv"
)

func (exp ExpPlus) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	_, err := exp.Left.Generate(gen)
	if err != nil {
		return InvalidExpKind, fmt.Errorf("failed to generate code for left side of plus exp: %w", err)
	}
	gen.push(rax)
	_, err = exp.Right.Generate(gen)
	if err != nil {
		return InvalidExpKind, fmt.Errorf("failed to generate code for right side of plus exp: %w", err)
	}
	gen.pop(rbx)
	gen.add(rax, rbx)
	return StackExp, nil
}

func (exp ExpMultiply) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	_, err := exp.Left.Generate(gen)
	if err != nil {
		return InvalidExpKind, fmt.Errorf("failed to generate code for left side of multiply exp: %w", err)
	}
	gen.push(rax)
	_, err = exp.Right.Generate(gen)
	if err != nil {
		return InvalidExpKind, fmt.Errorf("failed to generate code for right side of multiply exp: %w ", err)
	}
	gen.pop(rbx)
	gen.mult(rax, rbx)
	return StackExp, nil
}

func (exp ExpParentheses) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	_, err := exp.Inside.Generate(gen)
	if err != nil {
		return InvalidExpKind, fmt.Errorf("failed to generate code for inside of parentheses exp: %w", err)
	}
	return StackExp, nil
}

func (exp ExpNum) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	val := strconv.Itoa(exp.Value)
	gen.mov(rax, fmt.Sprintf("%s", val))
	return StackExp, nil
}

func (exp ExpIdentifier) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	kind, addr, err := gen.contexts.Get(exp.Name, gen.stackSize)
	if err != nil {
		return InvalidExpKind, fmt.Errorf("failed to evaluate identifer: %w", err)
	}
	gen.mov(rax, addr)

	translatedKind := InvalidExpKind
	if kind == assembly.ProcedureElem {
		translatedKind = ProcExp
	} else {
		translatedKind = StackExp
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
	if kind == StackExp {
		operations, pushes, err := gen.contexts.StackInsert(stmt.Identifier, rax, gen.stackSize)
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
	_, err := stmt.Expression.Generate(gen)
	if err != nil {
		return fmt.Errorf("failed to generate code for expression in println: %w", err)
	}
	gen.println(rax)
	return nil
}

func (stmt StmtReturn) Generate(gen *AssemblyGenerator) error {
	_, err := stmt.Expression.Generate(gen)
	if err != nil {
		return fmt.Errorf("failed to evaluate expression when returning: %w", err)
	}
	procedure := gen.contexts.GetTopProcedure()

	if procedure == nil {
		return fmt.Errorf("cant return outside a procedure")
	}

	gen.jmp(fmt.Sprintf("%send", procedure.Name))
	return nil
}

func (exp ExpFunction) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	newContext, operations, pushes := gen.contexts.NewContext(true, gen.stackSize)
	gen.AddOperations(operations)
	gen.stackSize += pushes
	gen.contexts.Push(newContext)

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
		return InvalidExpKind, fmt.Errorf("failed to generate function body: %w", err)
	}

	if len(exp.Args) != 0 {
		gen.stackSize -= len(exp.Args) + 1
	}

	pops, context := gen.contexts.Pop(gen.stackSize)
	context.Procedure.NumberOfArgs = len(exp.Args)
	context.Procedure.ReturnType = exp.ReturnType
	gen.stackSize -= pops
	gen.PushNamelessProcedure(context.Procedure)

	return ProcExp, nil
}

func (stmt FunctionCall) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	// Bind arguments
	for _, arg := range stmt.Arguments {
		arg.Generate(gen)
		gen.push(rax)
	}
	gen.stackSize -= len(stmt.Arguments)

	kind, actualName, err := gen.contexts.Get(stmt.Name, gen.stackSize)
	if kind != assembly.ProcedureElem || err != nil {
		return InvalidExpKind, fmt.Errorf("failed to call procedure with name: %s", stmt.Name)
	}

	var procedure *assembly.Procedure
	for i := range gen.AllProcedures {
		p := gen.AllProcedures[i]
		if p.Name == actualName {
			procedure = p
		}
	}
	if procedure == nil {
		return InvalidExpKind, fmt.Errorf("failed to find procedure with name: %s", actualName)
	}

	if procedure.NumberOfArgs != len(stmt.Arguments) {
		return InvalidExpKind, fmt.Errorf("mismatching number of arguments")
	}

	gen.mov(rcx, strconv.Itoa(gen.stackSize))
	gen.addOperation(fmt.Sprintf("call %s", actualName))

	// Pop bound arguments
	for range stmt.Arguments {
		gen.popWithoutDecreasingStackSize(rbx)
	}

	if err != nil {
		return InvalidExpKind, fmt.Errorf("failed to call function: %w", err)
	}

	if procedure.ReturnType == "empty" {
		return EmptyExp, nil
	}

	return StackExp, nil
}
