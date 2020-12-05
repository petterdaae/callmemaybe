package test

import (
	"lang/language"
	"strings"
	"testing"
)

func TestIdentifier(t *testing.T) {
	tokenizer := language.NewTokenizer(strings.NewReader("foo=bar"))
	first, _ := tokenizer.NextToken()
	second, _ := tokenizer.NextToken()
	third, _ := tokenizer.NextToken()

	if first != language.Identifier || second != language.Assign || third != language.Identifier {
		t.Error()
	}
}

func TestIdentifierStartingWithDigit(t *testing.T) {
	tokenizer := language.NewTokenizer(strings.NewReader("4sound"))
	first, _ := tokenizer.NextToken()

	if first == language.Identifier {
		t.Error()
	}

	tokenizer = language.NewTokenizer(strings.NewReader("proc1"))
	first, text := tokenizer.NextToken()
	if first != language.Identifier || text != "proc1" {
		t.Error()
	}
}

func TestPrintln(t *testing.T) {
	tokenizer := language.NewTokenizer(strings.NewReader("println 1"))
	first, _ := tokenizer.NextToken()
	if first != language.PrintLn {
		t.Error()
	}
}

func TestArrow(t *testing.T) {
	tokenizer := language.NewTokenizer(strings.NewReader("=>"))
	first, _ := tokenizer.NextToken()
	if first != language.Arrow {
		t.Error()
	}
}

func TestArrowWithSpace(t *testing.T) {
	tokenizer := language.NewTokenizer(strings.NewReader("= >"))
	first, _ := tokenizer.NextToken()
	if first == language.Arrow || first != language.Assign {
		t.Error()
	}
}