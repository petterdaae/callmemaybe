package parser

import (
	"lang/grammar"
	"reflect"
	"strings"
	"testing"
)

func parseExpected(t *testing.T, program string, expected grammar.Exp) {
	reader := strings.NewReader(program)
	parser := New(reader)
	actual := parser.Parse()
	if !reflect.DeepEqual(actual, expected) {
		t.Error()
	}
}

func TestSimplePlus(t *testing.T) {
	str := "1 + 2"
	expected := grammar.ExpPlus{
		Left:  grammar.ExpNum{Value: 1},
		Right: grammar.ExpNum{Value: 2},
	}
	parseExpected(t, str, expected)
}

func TestSimpleMultiply(t *testing.T) {
	str := "1 * 2"
	expected := grammar.ExpMultiply{
		Left:  grammar.ExpNum{Value: 1},
		Right: grammar.ExpNum{Value: 2},
	}
	parseExpected(t, str, expected)
}

func TestSimpleParentheses(t *testing.T) {
	str := "( 1 )"
	expected := grammar.ExpParentheses{
		Inside: grammar.ExpNum{Value: 1},
	}
	parseExpected(t, str, expected)
}

func TestParenthesesInPlusExpression(t *testing.T) {
	str := "( 1 + 2 ) + 3"
	expected := grammar.ExpPlus{
		Left: grammar.ExpParentheses{
			Inside: grammar.ExpPlus{
				Left:  grammar.ExpNum{1},
				Right: grammar.ExpNum{Value: 2},
			},
		},
		Right: grammar.ExpNum{Value: 3},
	}
	parseExpected(t, str, expected)
}

func TestBinaryExpressionWithTrailingParentheses(t *testing.T) {
	str := "1 * ( 2 * 3)"
	expected := grammar.ExpMultiply{
		Left: grammar.ExpNum{Value: 1},
		Right: grammar.ExpParentheses{
			Inside: grammar.ExpMultiply{
				Left:  grammar.ExpNum{Value: 2},
				Right: grammar.ExpNum{Value: 3},
			},
		},
	}
	parseExpected(t, str, expected)
}
