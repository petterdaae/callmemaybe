package grammar

type Stmt interface {
	Execute(ctx Context) (Context, error)
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
