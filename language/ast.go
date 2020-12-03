package language

type Exp interface {
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
	Identifier    string
	IdentifierExp Exp
	Inside        Exp
}

type ExpIdentifier struct {
	Name string
}

type ExpFunction struct {
	Args []Arg
	Body Stmt
}

type Arg struct {
	Identifier string
	Type string
}

type FunctionCall struct {
	Name string
	Arguments []Exp
}

type Stmt interface {
	Generate(output *AssemblyOutput) error
}

type StmtSeq struct {
	Statements []Stmt
}

type StmtAssign struct {
	Identifier string
	Expression Exp
}

type StmtPrintln struct {
	Expression Exp
}

type StmtReturn struct {
	Expression Exp
}

type ExprStmt interface {
	Generate(output *AssemblyOutput) error
}