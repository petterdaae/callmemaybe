package lexer

import (
	"reflect"
	"testing"
)

func expectedLexItems(kinds []LexItemKind, values []string) []LexItem {
	var result []LexItem
	for i := 0; i < len(kinds) && i < len(values); i++ {
		item := LexItem{
			kind:  kinds[i],
			value: []rune(values[i]),
		}
		result = append(result, item)
	}
	return result
}

func TestLexNumberSimple(t *testing.T) {
	lexer := New("123")
	LexNumber(&lexer)
	expected := expectedLexItems(
		[]LexItemKind{Number},
		[]string{"123"},
	)
	if !reflect.DeepEqual(lexer.lexItems, expected) {
		t.Error()
	}
	if lexer.currentIndex != 3 {
		t.Error()
	}
}

func TestLexNumberTrailingSpace(t *testing.T) {
	lexer := New("123   ")
	LexNumber(&lexer)
	expected := expectedLexItems(
		[]LexItemKind{Number},
		[]string{"123"},
	)
	if !reflect.DeepEqual(lexer.lexItems, expected) {
		t.Error()
	}
	if lexer.currentIndex != 3 {
		t.Error()
	}
}


func TestLexNumberFails(t *testing.T) {
	lexer := New(" 123")
	err := LexNumber(&lexer)
	if err == nil {
		t.Error()
	}
}

func TestLexNumberFailsOnEmpty(t *testing.T) {
	lexer := New("")
	err := LexNumber(&lexer)
	if err == nil {
		t.Error()
	}
}

func TestLexOperatorSimple(t *testing.T) {
	lexer := New("+")
	LexOperator(&lexer)
	expected := expectedLexItems(
		[]LexItemKind{Operator},
		[]string{"+"},
	)
	if !reflect.DeepEqual(lexer.lexItems, expected) {
		t.Error()
	}
	if lexer.currentIndex != 1 {
		t.Error()
	}
}

func TestLexOperatorFailsOnEmpty(t *testing.T) {
	lexer := New("")
	err := LexOperator(&lexer)
	if err == nil {
		t.Error()
	}
}


func TestLexOperatorTrailingSpace(t *testing.T) {
	lexer := New("*   ")
	LexOperator(&lexer)
	expected := expectedLexItems(
		[]LexItemKind{Operator},
		[]string{"*"},
	)
	if !reflect.DeepEqual(lexer.lexItems, expected) {
		t.Error()
	}
	if lexer.currentIndex != 1 {
		t.Error()
	}
}

func TestLexParenthesesSimple(t *testing.T) {
	lexer := New(")")
	LexParentheses(&lexer)
	expected := expectedLexItems(
		[]LexItemKind{Parentheses},
		[]string{")"},
	)
	if !reflect.DeepEqual(lexer.lexItems, expected) {
		t.Error()
	}
	if lexer.currentIndex != 1 {
		t.Error()
	}
}

func TestLexParenthesesFailsOnEmpty(t *testing.T) {
	lexer := New("")
	err := LexParentheses(&lexer)
	if err == nil {
		t.Error()
	}
}


func TestLexParenthesesTrailingSpace(t *testing.T) {
	lexer := New("(   ")
	LexParentheses(&lexer)
	expected := expectedLexItems(
		[]LexItemKind{Parentheses},
		[]string{"("},
	)
	if !reflect.DeepEqual(lexer.lexItems, expected) {
		t.Error()
	}
	if lexer.currentIndex != 1 {
		t.Error()
	}
}


func TestLexWhiteSpaceSimple(t *testing.T) {
	lexer := New(" \t\n")
	LexWhiteSpace(&lexer)
	if lexer.currentIndex != 3 {
		t.Error()
	}
}

func TestLexWhiteSpaceTrailingNumber(t *testing.T) {
	lexer := New("   1")
	LexWhiteSpace(&lexer)
	if lexer.currentIndex != 3 {
		t.Error()
	}
}


func TestLexWhiteSpaceFails(t *testing.T) {
	lexer := New("123 ")
	err := LexWhiteSpace(&lexer)
	if err == nil {
		t.Error()
	}
}

func TestLexWhiteSpaceFailsOnEmpty(t *testing.T) {
	lexer := New("")
	err := LexWhiteSpace(&lexer)
	if err == nil {
		t.Error()
	}
}
