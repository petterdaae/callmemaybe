package memorymodel

import (
	"callmemaybe/language/common"
)

type Context struct {
	members map[string]*ContextElement
}

type ContextElement struct {
	Kind                 common.ContextElementKind
	Name                 string
	StackSizeAfterPush   int
	FunctionNumberOfArgs int
	FunctionReturnKind   common.ContextElementKind
	FunctionArguments    []common.Arg
	ListElementKind      common.ContextElementKind
	ListSize             int
}

func EmptyContext() *Context {
	return &Context{
		members: make(map[string]*ContextElement),
	}
}

func NewContextElement(
	stackSizeAfterPush int,
	kind common.ContextElementKind,
	name string,
	numberOfArgs int,
	returnKind common.ContextElementKind,
	elementKind common.ContextElementKind,
	functionArguments []common.Arg,
	listSize int,
) *ContextElement {
	return &ContextElement{
		StackSizeAfterPush:   stackSizeAfterPush,
		Kind:                 kind,
		Name:                 name,
		FunctionNumberOfArgs: numberOfArgs,
		FunctionReturnKind:   returnKind,
		ListElementKind:      elementKind,
		FunctionArguments:    functionArguments,
		ListSize:             listSize,
	}
}
