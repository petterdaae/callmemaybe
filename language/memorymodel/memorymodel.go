package memorymodel

import (
	"callmemaybe/language/common"
)

type MemoryModel struct {
	CurrentStackSize int
	ContextStack     *ContextStack
}

func NewMemoryModel() *MemoryModel {
	return &MemoryModel{
		CurrentStackSize: 0,
		ContextStack:     NewContextStack(),
	}
}

func (mm *MemoryModel) PushNewContext(copyCurrentContext bool) {
	current := mm.ContextStack.Peek()
	newContext := EmptyContext()
	if copyCurrentContext {
		for k, v := range current.members {
			newContext.members[k] = v
		}
	}
	mm.ContextStack.Push(newContext)
}

func (mm *MemoryModel) PopCurrentContext() {
	mm.ContextStack.Pop()
}

func (mm *MemoryModel) AddNameToCurrentStackElement(name string, kind common.ContextElementKind, listElementKind common.ContextElementKind, listSize int) {
	currentContext := mm.ContextStack.Peek()
	currentContext.members[name] = NewContextElement(mm.CurrentStackSize, kind, "", 0, common.ContextElementKindInvalid, listElementKind, []common.Arg{}, listSize)
}

func (mm *MemoryModel) GetStackElement(name string) *ContextElement {
	value, ok := mm.ContextStack.Peek().members[name]
	if !ok {
		return nil
	}
	return value
}

func (mm *MemoryModel) GetProcedureElement(name string) *ContextElement {
	value, ok := mm.ContextStack.Peek().members[name]
	if !ok || value.Kind != common.ContextElementKindProcedure {
		return nil
	}
	return value
}

func (mm *MemoryModel) AddProcedureAlias(name string, alias string, numberOfArgs int, returnType common.ContextElementKind, args []common.Arg) {
	currentContext := mm.ContextStack.Peek()
	currentContext.members[alias] = NewContextElement(0, common.ContextElementKindProcedure, name, numberOfArgs, returnType, common.ContextElementKindInvalid, args, 0)
}

