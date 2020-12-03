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

func TestPrintln(t *testing.T) {
	tokenizer := language.NewTokenizer(strings.NewReader("println 1"))
	first, _ := tokenizer.NextToken()
	if first != language.PrintLn {
		t.Error()
	}
}