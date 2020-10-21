package grammar

type Exp interface{
	Evaluate(ctx Context) (int, error)
}

type ExpWithParent struct {
	Parent Exp
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

type ExpLet struct {
	Identifier string
	IdentifierExp Exp
	Inside Exp
}

type ExpIdentifier struct {
	Name string
}
