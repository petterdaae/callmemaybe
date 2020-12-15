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
	Struct
)

var passable = []RawType{Int, Char, Bool, List, Function, Struct}
var comparable = []RawType{Int, Char, Bool}

type Type struct {
	RawType               RawType
	ListElementType       *Type
	FunctionArgumentTypes []NamedType
	FunctionReturnType    *Type
	StructName            string
	StructMembers         []NamedType
}

type NamedType struct {
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

func (t Type) IsStorableOnStack() bool {
	return t.RawType != Invalid && t.RawType != Void
}

func (t Type) Equals(o Type) bool {
	if t.RawType != o.RawType {
		return false
	}

	if t.RawType == List {
		if !(*t.ListElementType).Equals(*o.ListElementType) {
			return false
		}
	}

	if t.RawType == Function {
		if !(*t.FunctionReturnType).Equals(*o.FunctionReturnType) {
			return false
		}

		if len(t.FunctionArgumentTypes) != len(o.FunctionArgumentTypes) {
			return false
		}

		for i := 0; i < len(t.FunctionArgumentTypes); i++ {
			tArg := t.FunctionArgumentTypes[i]
			oArg := o.FunctionArgumentTypes[i]
			if !tArg.Type.Equals(oArg.Type) {
				return false
			}
		}
	}

	if t.RawType == Struct {
		if t.StructName != o.StructName {
			return false
		}
	}

	return true
}

func NewInt() Type {
	return Type{
		RawType: Int,
	}
}

func NewChar() Type {
	return Type{
		RawType: Char,
	}
}

func NewBool() Type {
	return Type{
		RawType: Bool,
	}
}

func NewInvalid() Type {
	return Type{
		RawType: Invalid,
	}
}

func contains(elem RawType, collection []RawType) bool {
	for _, item := range collection {
		if item == elem {
			return true
		}
	}
	return false
}
