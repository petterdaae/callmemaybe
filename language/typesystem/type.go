package typesystem

type RawType int

const (
	Invalid RawType = iota
	Void
	Int
	Char
	Bool
	List
	Function
)

var passable = []RawType{Int, Char, Bool, List, Function}
var comparable = []RawType{Int, Char, Bool}

type Type struct {
	RawType               RawType
	ListElementType       *Type
	ListSize              int
	FunctionArgumentTypes []FunctionArgument
	FunctionReturnType    *Type
}

type FunctionArgument struct {
	Name string
	Type Type
}

func (t Type) IsPassable() bool {
	return contains(t.RawType, passable)
}

func (t Type) IsAlgebraic() bool {
	return t.RawType == Int
}

func (t Type) IsComparable() bool {
	return contains(t.RawType, comparable)
}

func contains(elem RawType, collection []RawType) bool {
	for _, item := range collection {
		if item == elem {
			return true
		}
	}
	return false
}
