package language

import "fmt"

func (exp ExpPlus) Evaluate(ctx Context) (int, error) {
	left, leftErr := exp.Left.Evaluate(ctx)
	if leftErr != nil {
		return left, fmt.Errorf("failed to evaluate left side of plus expression: %w", leftErr)
	}
	right, rightErr := exp.Right.Evaluate(ctx)
	if rightErr != nil {
		return right, fmt.Errorf("failed to evaluate right side of plus expression: %w", rightErr)
	}
	return left + right, nil
}

func (exp ExpMultiply) Evaluate(ctx Context) (int, error) {
	left, leftErr := exp.Left.Evaluate(ctx)
	if leftErr != nil {
		return left, fmt.Errorf("failed to evaluate left side of multiply expression: %w", leftErr)
	}
	right, rightErr := exp.Right.Evaluate(ctx)
	if rightErr != nil {
		return right, fmt.Errorf("failed to evaluate right side of multiply expression: %w", rightErr)
	}
	return left * right, nil
}

func (exp ExpParentheses) Evaluate(ctx Context) (int, error) {
	return exp.Inside.Evaluate(ctx)
}

func (exp ExpNum) Evaluate(ctx Context) (int, error) {
	return exp.Value, nil
}

func (exp ExpIdentifier) Evaluate(ctx Context) (int, error) {
	value, err := ctx.Get(exp.Name)
	return value, err
}

func (exp ExpLet) Evaluate(ctx Context) (int, error) {
	identifierValue, err := exp.IdentifierExp.Evaluate(ctx)
	if err != nil {
		return identifierValue, fmt.Errorf("failed to evaluate first expression in let: %w", err)
	}
	newContext := ctx.CopyWith(exp.Identifier, identifierValue)
	result, err := exp.Inside.Evaluate(newContext)
	if err != nil {
		return result, fmt.Errorf("failed to evaluate last expression in let: %w", err)
	}
	return result, nil
}
