package language

type Exp interface{
	Generate(output *AssemblyOutput) error
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
