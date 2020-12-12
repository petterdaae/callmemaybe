package language

import (
	"fmt"
	"lang/language/common"
)

type GenerateResult struct {
	Kind            common.ContextElementKind
	ListElementKind common.ContextElementKind
	ProcedureName   string
	Error           error
}

func (gr GenerateResult) IsError() bool {
	return gr.Error != nil
}

func (gr GenerateResult) WrapError(message string) GenerateResult {
	return GenerateResult{
		Kind: common.ContextElementKindInvalid,
		Error: fmt.Errorf("%s: %w", message, gr.Error),
	}
}

func (gr GenerateResult) Copy() GenerateResult {
	return GenerateResult{
		Kind:            gr.Kind,
		ListElementKind: gr.ListElementKind,
		ProcedureName:   gr.ProcedureName,
		Error:           gr.Error,
	}
}

func NumberResult() GenerateResult {
	return GenerateResult{
		Kind: common.ContextElementKindNumber,
	}
}

func BoolResult() GenerateResult {
	return GenerateResult{
		Kind: common.ContextElementKindBoolean,
	}
}

func CharResult() GenerateResult {
	return GenerateResult{
		Kind: common.ContextElementKindChar,
	}
}

func ErrorResult(error error) GenerateResult {
	return GenerateResult{
		Kind:  common.ContextElementKindChar,
		Error: error,
	}
}

func ProcedureResult(name string) GenerateResult {
	return GenerateResult{
		Kind:          common.ContextElementKindProcedure,
		ProcedureName: name,
	}
}

func ListResult(elements common.ContextElementKind) GenerateResult {
	return GenerateResult{
		ListElementKind: elements,
		Kind:            common.ContextElementKindListReference,
	}
}

func CustomResult(kind common.ContextElementKind, listElementKind common.ContextElementKind, name string, err error) GenerateResult {
	return GenerateResult{
		Kind:            kind,
		ListElementKind: listElementKind,
		ProcedureName:   name,
		Error:           err,
	}
}
