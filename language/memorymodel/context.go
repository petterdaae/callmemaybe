package memorymodel

import (
	"callmemaybe/language/typesystem"
)

type Context struct {
	members map[string]*ContextElement
	structTypes map[string]typesystem.Type
}

type ContextElement struct {
	Type                 typesystem.Type
	Name                 string
	StackSizeAfterPush   int
}

func EmptyContext() *Context {
	return &Context{
		members: make(map[string]*ContextElement),
		structTypes: make(map[string]typesystem.Type),
	}
}

func NewContextElement(
	_type typesystem.Type,
	stackSizeAfterPush int,
	name string,
) *ContextElement {
	return &ContextElement{
		Type: _type,
		StackSizeAfterPush:   stackSizeAfterPush,
		Name:                 name,
	}
}
