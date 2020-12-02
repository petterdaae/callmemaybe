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

func TestLet(t *testing.T) {
	tokenizer := language.NewTokenizer(strings.NewReader("let+in"))
	first, _ := tokenizer.NextToken()
	tokenizer.NextToken()
	third, _ := tokenizer.NextToken()

	if first != language.Let || third != language.In {
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