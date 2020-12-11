package language

import "lang/language/memorymodel"

type GenerationResult struct {
	Kind memorymodel.ContextElementKind
	ListElementKind memorymodel.ContextElementKind
	ProcedureName string
	Error error
}

func NumberKind() GenerationResult {
	return GenerationResult{
		Kind: memorymodel.ContextElementKindNumber,
	}
}

func BoolKind() GenerationResult {
	return GenerationResult{
		Kind: memorymodel.ContextElementKindBoolean,
	}
}

func CharKind() GenerationResult {
	return GenerationResult{
		Kind: memorymodel.ContextElementKindChar,
	}
}

func ErrorKind(error error) GenerationResult {
	return GenerationResult{
		Kind: memorymodel.ContextElementKindChar,
		Error: error,
	}
}

func ProcedureKind(name string) GenerationResult {
	return GenerationResult{
		Kind: memorymodel.ContextElementKindProcedure,
		ProcedureName: name,
	}
}

func ListKind(elements memorymodel.ContextElementKind) GenerationResult {
	return GenerationResult{
		ListElementKind: elements,
		Kind: memorymodel.ContextElementKindListReference,
	}
}

func CustomKind(kind memorymodel.ContextElementKind, listElementKind memorymodel.ContextElementKind, name string, err error) GenerationResult {
	return GenerationResult{
		Kind: kind,
		Error: err,
		ProcedureName: name,
		ListElementKind: listElementKind,
	}
}