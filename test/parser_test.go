package test

import (
	"lang/language"
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

func testSuccessfulParser(t *testing.T, program string, expected int) {
	parser := language.NewParser(strings.NewReader(program))
	exp, err := parser.ParseExp()
	if err != nil {
		t.Error()
	}
	val, err := exp.Evaluate(language.NewContext())
	if val != expected || err != nil {
		t.Error()
	}
}

func testFailingParser(t *testing.T, program string) {
	parser := language.NewParser(strings.NewReader(program))
	exp, parseErr := parser.ParseExp()
	if parseErr != nil {
		return
	}
	_, evalErr := exp.Evaluate(language.NewContext())
	if evalErr != nil {
		return
	}
	t.Error()
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

func TestSimpleLeftAssociativity(t *testing.T) {
	parser := language.NewParser(strings.NewReader("(1 + 4) * 2 + 3 * 5"))
	expr, err := parser.ParseExp()
	if err != nil {
		t.Error()
	}
	val, _ := expr.Evaluate(language.NewContext())
	if val != 65 {
		t.Error()
	}
}

func TestSimpleLet(t *testing.T) {
	testSuccessfulParser(t, "let a = 2 in a", 2)
}

func TestNestedLet(t *testing.T) {
	testSuccessfulParser(t, "let a = 2 in let a = 3 in 3", 3)
}

func TestVariousLet(t *testing.T) {
	testSuccessfulParser(t, "let a = 5 in 2 + a", 7)
	testSuccessfulParser(t, "7 * (let b = 2 in b + 5)", 49)
	testSuccessfulParser(t, "(let a = 2 in let b = a + 3 in b + 2)", 7)
	testSuccessfulParser(t, "(let a = 1 + 2 + 3 in a + 4) + 5", 15)
	testSuccessfulParser(t, "let foo = (let bar = 3 in bar * 2) in foo", 6)
}

func TestVariousFailingLet(t *testing.T) {
	testFailingParser(t, "let a")
	testFailingParser(t, "let x = 5 in y")
	testFailingParser(t, "let x = 5 in")
	testFailingParser(t, "let x3 = 5 in x3")
}
