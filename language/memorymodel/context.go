package memorymodel

type ContextElementKind int

const (
	ContextElementKindNumber ContextElementKind = iota
	ContextElementKindBoolean
	ContextElementKindInvalid
	ContextElementKindProcedure
	ContextElementKindEmpty
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
}

func EmptyContext() *Context {
	return &Context{
		members: make(map[string]*ContextElement),
	}
}

func NewContextElement(stackSizeAfterPush int, kind ContextElementKind, name string, numberOfArgs int, returnKind ContextElementKind) *ContextElement {
	return &ContextElement{
		StackSizeAfterPush: stackSizeAfterPush,
		Kind:               kind,
		Name:               name,
		NumberOfArgs:       numberOfArgs,
		ReturnKind:         returnKind,
	}
}

func IsStackKind(kind ContextElementKind) bool {
	return kind != ContextElementKindInvalid && kind != ContextElementKindProcedure
}

func GetKindFromType(_type string) ContextElementKind {
	switch _type {
	case "bool":
		return ContextElementKindBoolean
	case "int":
		return ContextElementKindNumber
	default:
		return ContextElementKindInvalid
	}
}
