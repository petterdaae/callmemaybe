package grammar

type Exp interface{
	Evaluate() int
}

type ExpPlus struct {
	Left  Exp
	Right Exp
}

type ExpMultiply struct {
	Left  Exp
	Right Exp
}

type ExpParentheses struct {
	Inside Exp
}

type ExpNum struct {
	Value int
}

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
