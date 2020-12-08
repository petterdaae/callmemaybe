package memorymodel

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
	newContext := EmptyContext()
	if copyCurrentContext {
		for k, v := range newContext.members {
			newContext.members[k] = v
		}
	}
	mm.ContextStack.Push(newContext)
}

func (mm *MemoryModel) PopCurrentContext() []*ContextElement {
	popped := mm.ContextStack.Peek()
	mm.ContextStack.Pop()
	current := mm.ContextStack.Peek()
	var unavailableStackElements []*ContextElement
	for k, v := range popped.members {
		_, ok := current.members[k]
		if !ok {
			unavailableStackElements = append(unavailableStackElements, v)
		}
	}
	return unavailableStackElements
}

func (mm *MemoryModel) AddNameToCurrentStackElement(name string, kind ContextElementKind) {
	currentContext := mm.ContextStack.Peek()
	currentContext.members[name] = NewContextElement(mm.CurrentStackSize, kind, "", 0, ContextElementKindInvalid)
}

func (mm *MemoryModel) AddNameToStackElement(name string, kind ContextElementKind, stackSizeWhenPushed int) {
	currentContext := mm.ContextStack.Peek()
	currentContext.members[name] = NewContextElement(stackSizeWhenPushed, kind, "", 0, ContextElementKindInvalid)
}

func (mm *MemoryModel) GetStackElement(name string) *ContextElement {
	value, ok := mm.ContextStack.Peek().members[name]
	if !ok || !IsStackKind(value.Kind) {
		return nil
	}
	return value
}

func (mm *MemoryModel) GetProcedureElement(name string) *ContextElement {
	value, ok := mm.ContextStack.Peek().members[name]
	if !ok || value.Kind != ContextElementKindProcedure {
		return nil
	}
	return value
}

func (mm *MemoryModel) AddProcedureAlias(name string, alias string, numberOfArgs int, returnType ContextElementKind) {
	currentContext := mm.ContextStack.Peek()
	currentContext.members[alias] = NewContextElement(0, ContextElementKindProcedure, name, numberOfArgs, returnType)
}

