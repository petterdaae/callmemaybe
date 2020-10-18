package grammar

func (exp ExpPlus) Evaluate() int {
	left := exp.Left.Evaluate()
	right := exp.Right.Evaluate()
	return left + right
}

func (exp ExpMultiply) Evaluate() int {
	left := exp.Left.Evaluate()
	right := exp.Right.Evaluate()
	return left * right
}

func (exp ExpParentheses) Evaluate() int {
	return exp.Inside.Evaluate()
}

func (exp ExpNum) Evaluate() int {
	return exp.Value
}
