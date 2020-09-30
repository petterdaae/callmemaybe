package grammar

type Expression interface{
	Evaluate() int
}

type ExpPlus struct {
	Left  Expression
	Right Expression
}

type ExpMultiply struct {
	Left  Expression
	Right Expression
}

type ExpParentheses struct {
	Inside Expression
}

func (exp *ExpPlus) Evaluate() int {
	left := exp.Left.Evaluate()
	right := exp.Right.Evaluate()
	return left + right
}

func (exp *ExpMultiply) Evaluate() int {
	left := exp.Left.Evaluate()
	right := exp.Right.Evaluate()
	return left * right
}

func (exp *ExpParentheses) Evaluate() int {
	return exp.Inside.Evaluate()
}