package language

import (
	"fmt"
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
	return InvalidExpKind, nil
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
	return InvalidExpKind, nil
}

func (exp ExpParentheses) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	_, err := exp.Inside.Generate(gen)
	if err != nil {
		return InvalidExpKind, fmt.Errorf("failed to generate code for inside of parentheses exp: %w", err)
	}
	return InvalidExpKind, nil
}

func (exp ExpNum) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	val := strconv.Itoa(exp.Value)
	gen.move(rax, fmt.Sprintf("%s", val))
	return InvalidExpKind, nil
}

func (exp ExpIdentifier) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	/*
	identifierStackPos, ok := gen.Identifiers[exp.Name]
	if !ok {
		return fmt.Errorf("uknown identifier: %s", exp.Name)
	}
	diff := (gen.StackSize - identifierStackPos) * 8
	if diff < 0 {
		return fmt.Errorf("negative stack position for identifier: %s", exp.Name)
	}
	identifierAddr := fmt.Sprintf("[%s+%d]", rsp, diff)
	*/
	_, addr, err := gen.get(exp.Name)
	if err != nil {
		return InvalidExpKind, fmt.Errorf("failed to evaluate identifer: %w", err)
	}
	gen.move(rax, addr)
	return InvalidExpKind, nil
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
	_, err := stmt.Expression.Generate(gen)
	if err != nil {
		return fmt.Errorf("failed to generate code for expression in assign statement: %w", err)
	}
	gen.push(rax)
	gen.pushToStack(stmt.Identifier)
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
	// TODO : implement
	return nil
}

func (stmt ExpFunction) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	// TODO : implement
	return InvalidExpKind, nil
}

func (stmt FunctionCall) Generate(gen *AssemblyGenerator) (ExpKind, error) {
	// TODO : implement
	return InvalidExpKind, nil
}
