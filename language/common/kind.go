package common

type ContextElementKind int

const (
	ContextElementKindInvalid ContextElementKind = iota
	ContextElementKindNumber
	ContextElementKindBoolean
	ContextElementKindChar
	ContextElementKindProcedure
	ContextElementKindEmpty
	ContextElementKindListReference
)

func (c ContextElementKind) IsPassable() bool {
	switch c {
	case ContextElementKindNumber:
		return true
	case ContextElementKindChar:
		return true
	case ContextElementKindBoolean:
		return true
	default:
		return false
	}
}

func (c ContextElementKind) IsAlgebraic() bool {
	switch c {
	case ContextElementKindNumber:
		return true
	default:
		return false
	}
}

func (c ContextElementKind) IsComparable() bool {
	switch c {
	case ContextElementKindNumber:
		return true
	case ContextElementKindChar:
		return true
	default:
		return false
	}
}

func (c ContextElementKind) IsStoredOnStack() bool {
	switch c {
	case ContextElementKindNumber:
		return true
	case ContextElementKindChar:
		return true
	case ContextElementKindBoolean:
		return true
	case ContextElementKindListReference:
		return true
	default:
		return false
	}
}
