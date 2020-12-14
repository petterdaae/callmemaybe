package memorymodel

import (
	"callmemaybe/language/typesystem"
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

func (mm *MemoryModel) AddNameToCurrentStackElement(name string, _type typesystem.Type) {
	currentContext := mm.ContextStack.Peek()
	currentContext.members[name] = NewContextElement(_type, mm.CurrentStackSize, name)
}

func (mm *MemoryModel) Update(name string, _type typesystem.Type) {
	currentContext := mm.ContextStack.Peek()
	member, _ := currentContext.members[name]
	currentContext.members[name] = NewContextElement(_type, member.StackSizeAfterPush, name)
}

func (mm *MemoryModel) Contains(name string) bool {
	currentContext := mm.ContextStack.Peek()
	_, ok := currentContext.members[name]
	return ok
}

func (mm *MemoryModel) GetStackElement(name string) *ContextElement {
	value, ok := mm.ContextStack.Peek().members[name]
	if !ok {
		return nil
	}
	return value
}
