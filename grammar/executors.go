package grammar

import "fmt"

func (stmt *StmtSeq) Execute(ctx Context) (Context, error) {
	var err error
	currentCtx := ctx
	for i := range stmt.Statements {
		current := stmt.Statements[i]
		currentCtx, err = current.Execute(currentCtx)
		if err != nil {
			return NewContext(), fmt.Errorf("failed to execute statement: %w", err)
		}
	}
	return currentCtx, nil
}

func (stmt *StmtAssign) Execute(ctx Context) (Context, error) {
	val, err := stmt.Expression.Evaluate(ctx)
	if err != nil {
		return NewContext(), fmt.Errorf("failed to evaluate expression in assign statement: %w", err)
	}
	newContext := ctx.CopyWith(stmt.Identifier, val)
	return newContext, nil
}

func (stmt *StmtPrintln) Execute(ctx Context) (Context, error) {
	val, err := stmt.Expression.Evaluate(ctx)
	if err != nil {
		return NewContext(), fmt.Errorf("failed to evaluate expression in println statement: %w", err)
	}
	println(val)
	return ctx, nil
}
