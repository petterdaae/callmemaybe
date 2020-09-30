package parser

import (
	"lang/lexer"
	"testing"
)

func testProgramEvaluatesTo(text string, result int, t *testing.T) {
	program, _ := lexer.Lex(text)
	exp, _ := Parse(program)
	if exp.Evaluate() != result {
		t.Error()
	}
}

func TestParserCase1(t *testing.T) {
	testProgramEvaluatesTo("1 + 2", 3, t)
}

func TestParserCase2(t *testing.T) {
	testProgramEvaluatesTo("(1 + 2)", 3, t)
}

func TestParserCase3(t *testing.T) {
	testProgramEvaluatesTo("1 + 2 + 3", 6, t)
}
