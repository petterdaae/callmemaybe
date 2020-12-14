package test

import (
	"callmemaybe/language"
	"callmemaybe/language/typesystem"
	"reflect"
	"strings"
	"testing"
)

func parseExpected(t *testing.T, program string, expected language.Exp) {
	reader := strings.NewReader(program)
	parser := language.NewParser(reader)
	actual, err := parser.ParseExp()

	if err != nil {
		t.Error()
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Error()
	}
}

func parseExpectedStmt(t *testing.T, program string, expected language.Stmt) {
	reader := strings.NewReader(program)
	parser := language.NewParser(reader)
	actual, err := parser.Parse()

	if err != nil {
		t.Error()
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Error()
	}
}

func TestSimplePlus(t *testing.T) {
	str := "1 + 2"
	expected := language.ExpPlus{
		Left:  language.ExpNum{Value: 1},
		Right: language.ExpNum{Value: 2},
	}
	parseExpected(t, str, expected)
}

func TestSimpleMultiply(t *testing.T) {
	str := "1 * 2"
	expected := language.ExpMultiply{
		Left:  language.ExpNum{Value: 1},
		Right: language.ExpNum{Value: 2},
	}
	parseExpected(t, str, expected)
}

func TestSimpleParentheses(t *testing.T) {
	str := "( 1 )"
	expected := language.ExpParentheses{
		Inside: language.ExpNum{Value: 1},
	}
	parseExpected(t, str, expected)
}

func TestParenthesesInPlusExpression(t *testing.T) {
	str := "( 1 + 2 ) + 3"
	expected := language.ExpPlus{
		Left: language.ExpParentheses{
			Inside: language.ExpPlus{
				Left:  language.ExpNum{1},
				Right: language.ExpNum{Value: 2},
			},
		},
		Right: language.ExpNum{Value: 3},
	}
	parseExpected(t, str, expected)
}

func TestBinaryExpressionWithTrailingParentheses(t *testing.T) {
	str := "1 * ( 2 * 3)"
	expected := language.ExpMultiply{
		Left: language.ExpNum{Value: 1},
		Right: language.ExpParentheses{
			Inside: language.ExpMultiply{
				Left:  language.ExpNum{Value: 2},
				Right: language.ExpNum{Value: 3},
			},
		},
	}
	parseExpected(t, str, expected)
}

func TestEmptyProgram(t *testing.T) {
	parser := language.NewParser(strings.NewReader(""))
	_, err := parser.ParseExp()
	if err == nil {
		t.Error()
	}
}

func TestMissingClosingParentheses(t *testing.T) {
	parser := language.NewParser(strings.NewReader("(1 + 2"))
	_, err := parser.ParseExp()
	if err == nil {
		t.Error()
	}
}

func TestEmptyParentheses(t *testing.T) {
	parser := language.NewParser(strings.NewReader("()"))
	_, err := parser.ParseExp()
	if err == nil {
		t.Error()
	}
}

func TestFunctionAssign(t *testing.T) {
	str := "f = | a int | { return a }"
	expected := language.StmtSeq{
		Statements: []language.Stmt{
			language.StmtAssign{
				Identifier: "f",
				Expression: language.ExpFunction{
					Body: language.StmtSeq{
						Statements: []language.Stmt{
							language.StmtReturn{
								Expression: language.ExpIdentifier{
									Name: "a",
								},
							},
						},
					},
					Type: typesystem.Type{
						RawType: typesystem.Function,
						FunctionArgumentTypes: []typesystem.FunctionArgument{typesystem.FunctionArgument{
							Name: "a",
							Type: typesystem.NewInt(),
						}},
						FunctionReturnType: &typesystem.Type{
							RawType:               typesystem.Void,
						},
					},
				},
			},
		},
	}
	parseExpectedStmt(t, str, expected)
}

func TestFunctionAssignWithType(t *testing.T) {
	str := "f = | a int | int { return a }"
	intType := typesystem.NewInt()
	expected := language.StmtSeq{
		Statements: []language.Stmt{
			language.StmtAssign{
				Identifier: "f",
				Expression: language.ExpFunction{
					Body: language.StmtSeq{
						Statements: []language.Stmt{
							language.StmtReturn{
								Expression: language.ExpIdentifier{
									Name: "a",
								},
							},
						},
					},
					Type: typesystem.Type{
						RawType: typesystem.Function,
						FunctionArgumentTypes: []typesystem.FunctionArgument{typesystem.FunctionArgument{
							Name: "a",
							Type: typesystem.NewInt(),
						}},
						FunctionReturnType: &intType,
					},
				},
			},
		},
	}
	parseExpectedStmt(t, str, expected)
}

func TestSimpleFunctionCall(t *testing.T) {
	str := "a = call func with 1, 2, 3"
	expected := language.StmtSeq{
		Statements: []language.Stmt{
			language.StmtAssign{
				Identifier: "a",
				Expression: language.FunctionCall{
					Exp: language.ExpIdentifier{Name: "func"},
					Arguments: []language.Exp{
						language.ExpNum{Value: 1},
						language.ExpNum{Value: 2},
						language.ExpNum{Value: 3},
					},
				},
			},
		},
	}
	parseExpectedStmt(t, str, expected)
}

func TestSimpleFunctionCallWithoutArgs(t *testing.T) {
	str := "a = call func"
	expected := language.StmtSeq{
		Statements: []language.Stmt{
			language.StmtAssign{
				Identifier: "a",
				Expression: language.FunctionCall{
					Exp:       language.ExpIdentifier{Name: "func"},
					Arguments: nil,
				},
			},
		},
	}
	parseExpectedStmt(t, str, expected)
}
