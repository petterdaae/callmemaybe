package language

type ExpKind int

const (
	StackExp ExpKind = iota
	ProcExp
	InvalidExpKind
)

type Exp interface {
	Generate(gen AssemblyGenerator) (ExpKind, error)
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
	ReturnType string
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
	Generate(gen AssemblyGenerator) error
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
