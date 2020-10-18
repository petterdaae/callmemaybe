package parser

import (
	"lang/tokenizer"
	"testing"
)

func testProgramEvaluatesTo(text string, result int, t *testing.T) {
	program, _ := tokenizer.Lex(text)
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

func TestParserCase4(t *testing.T) {
	testProgramEvaluatesTo("(((1*2)))", 2, t)
}

func TestParserCase5(t *testing.T) {
	testProgramEvaluatesTo("(1+2)+(3+3+3)+(1*2*3*4)", 36, t)
}

func TestParserCase6(t *testing.T) {
	testProgramEvaluatesTo("(1+2)+(3+4)", 10, t)
}