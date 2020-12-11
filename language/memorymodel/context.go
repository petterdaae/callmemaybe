package memorymodel

type ContextElementKind int

const (
	ContextElementKindNumber ContextElementKind = iota
	ContextElementKindBoolean
	ContextElementKindChar
	ContextElementKindInvalid
	ContextElementKindProcedure
	ContextElementKindEmpty
	ContextElementKindListReference
)

type Context struct {
	members map[string]*ContextElement
}

type ContextElement struct {
	StackSizeAfterPush int
	Kind               ContextElementKind
	Name               string
	NumberOfArgs       int
	ReturnKind         ContextElementKind
	ListElementKind    ContextElementKind
}

func (ce *ContextElement) Copy() *ContextElement {
	return &ContextElement{
		StackSizeAfterPush: ce.StackSizeAfterPush,
		Kind:               ce.Kind,
		Name:               ce.Name,
		NumberOfArgs:       ce.NumberOfArgs,
		ReturnKind:         ce.ReturnKind,
		ListElementKind:    ce.ListElementKind,
	}
}

func EmptyContext() *Context {
	return &Context{
		members: make(map[string]*ContextElement),
	}
}

func NewContextElement(stackSizeAfterPush int, kind ContextElementKind, name string, numberOfArgs int, returnKind ContextElementKind, elementKind ContextElementKind) *ContextElement {
	return &ContextElement{
		StackSizeAfterPush: stackSizeAfterPush,
		Kind:               kind,
		Name:               name,
		NumberOfArgs:       numberOfArgs,
		ReturnKind:         returnKind,
		ListElementKind:    elementKind,
	}
}

func IsIntOrBool(kind ContextElementKind) bool {
	return kind != ContextElementKindInvalid && kind != ContextElementKindProcedure
}

func GetKindFromType(_type string) ContextElementKind {
	switch _type {
	case "bool":
		return ContextElementKindBoolean
	case "int":
		return ContextElementKindNumber
	case "char":
		return ContextElementKindChar
	default:
		return ContextElementKindInvalid
	}
}
