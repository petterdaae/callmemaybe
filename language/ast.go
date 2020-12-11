package language

import (
	"lang/language/assemblyoutput"
	"lang/language/memorymodel"
)

type Exp interface {
	Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) GenerationResult
}

type ExpBop interface {
	LeftExp() Exp
	RightExp() Exp
}

type ExpPlus struct {
	Left  Exp
	Right Exp
}

func (exp ExpPlus) LeftExp() Exp {
	return exp.Left
}

func (exp ExpPlus) RightExp() Exp {
	return exp.Right
}

type ExpMinus struct {
	Left  Exp
	Right Exp
}

func (exp ExpMinus) LeftExp() Exp {
	return exp.Left
}

func (exp ExpMinus) RightExp() Exp {
	return exp.Right
}

type ExpDivide struct {
	Left  Exp
	Right Exp
}

func (exp ExpDivide) LeftExp() Exp {
	return exp.Left
}

func (exp ExpDivide) RightExp() Exp {
	return exp.Right
}

type ExpModulo struct {
	Left  Exp
	Right Exp
}

func (exp ExpModulo) LeftExp() Exp {
	return exp.Left
}

func (exp ExpModulo) RightExp() Exp {
	return exp.Right
}

type ExpNegative struct {
	Inside Exp
}

type ExpMultiply struct {
	Left  Exp
	Right Exp
}

func (exp ExpMultiply) LeftExp() Exp {
	return exp.Left
}

func (exp ExpMultiply) RightExp() Exp {
	return exp.Right
}

type ExpLess struct {
	Left  Exp
	Right Exp
}

func (exp ExpLess) LeftExp() Exp {
	return exp.Left
}

func (exp ExpLess) RightExp() Exp {
	return exp.Right
}

type ExpGreater struct {
	Left  Exp
	Right Exp
}

func (exp ExpGreater) LeftExp() Exp {
	return exp.Left
}

func (exp ExpGreater) RightExp() Exp {
	return exp.Right
}

type ExpEquals struct {
	Left  Exp
	Right Exp
}

func (exp ExpEquals) LeftExp() Exp {
	return exp.Left
}

func (exp ExpEquals) RightExp() Exp {
	return exp.Right
}

type ExpParentheses struct {
	Inside Exp
}

type ExpNum struct {
	Value int
}

type ExpChar struct {
	Value string
}

type ExpBool struct {
	Value bool
}

type ExpIdentifier struct {
	Name string
}

type ExpList struct {
	Elements []Exp
	Type     memorymodel.ContextElementKind
	Size     int
}

type ExpGetFromList struct {
	List  Exp
	Index Exp
}

type StmtAppendToList struct {
	List       Exp
	NewElement Exp
}

type ExpFunction struct {
	Recurse    string
	Args       []Arg
	Body       Stmt
	ReturnType memorymodel.ContextElementKind
}

type Arg struct {
	Identifier string
	Type       memorymodel.ContextElementKind
}

type FunctionCall struct {
	Name      string
	Arguments []Exp
}

type Stmt interface {
	Generate(ao *assemblyoutput.AssemblyOutput, mm *memorymodel.MemoryModel) error
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
