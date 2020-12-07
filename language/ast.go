package language

type ExpKind int

const (
	StackExp ExpKind = iota
	EmptyExp
	ProcExp
	InvalidExpKind
)

type Exp interface {
	Generate(gen *AssemblyGenerator) (ExpKind, error)
}

type ExpPlus struct {
	Left  Exp
	Right Exp
}

type ExpMultiply struct {
	Left  Exp
	Right Exp
}

type ExpLess struct {
	Left  Exp
	Right Exp
}

type ExpGreater struct {
	Left  Exp
	Right Exp
}

type ExpEquals struct {
	Left  Exp
	Right Exp
}

type ExpParentheses struct {
	Inside Exp
}

type ExpNum struct {
	Value int
}

type ExpBool struct {
	Value bool
}

type ExpIdentifier struct {
	Name string
}

type ExpFunction struct {
	Args       []Arg
	Body       Stmt
	ReturnType string
}

type Arg struct {
	Identifier string
	Type       string
}

type FunctionCall struct {
	Name      string
	Arguments []Exp
}

type Stmt interface {
	Generate(gen *AssemblyGenerator) error
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

type StmtIf struct {
	Expression Exp
	Body       Stmt
}
