package tokenizer

import (
	"reflect"
	"testing"
)

func expectedLexItems(kinds []Token, values []string) []LexItem {
	var result []LexItem
	for i := 0; i < len(kinds) && i < len(values); i++ {
		item := LexItem{
			Kind:  kinds[i],
			value: []rune(values[i]),
		}
		result = append(result, item)
	}
	return result
}

func TestLexNumberSimple(t *testing.T) {
	lexer := new("123")
	lexNumber(&lexer)
	expected := expectedLexItems(
		[]Token{Number},
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
	lexer := new("123   ")
	lexNumber(&lexer)
	expected := expectedLexItems(
		[]Token{Number},
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
	lexer := new(" 123")
	err := lexNumber(&lexer)
	if err == nil {
		t.Error()
	}
}

func TestLexNumberFailsOnEmpty(t *testing.T) {
	lexer := new("")
	err := lexNumber(&lexer)
	if err == nil {
		t.Error()
	}
}

func TestLexOperatorSimple(t *testing.T) {
	lexer := new("+")
	lexOperator(&lexer)
	expected := expectedLexItems(
		[]Token{Operator},
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
	lexer := new("")
	err := lexOperator(&lexer)
	if err == nil {
		t.Error()
	}
}


func TestLexOperatorTrailingSpace(t *testing.T) {
	lexer := new("*   ")
	lexOperator(&lexer)
	expected := expectedLexItems(
		[]Token{Operator},
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
	lexer := new(")")
	lexParentheses(&lexer)
	expected := expectedLexItems(
		[]Token{Parentheses},
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
	lexer := new("")
	err := lexParentheses(&lexer)
	if err == nil {
		t.Error()
	}
}


func TestLexParenthesesTrailingSpace(t *testing.T) {
	lexer := new("(   ")
	lexParentheses(&lexer)
	expected := expectedLexItems(
		[]Token{Parentheses},
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
	lexer := new(" \t\n")
	lexWhiteSpace(&lexer)
	if lexer.currentIndex != 3 {
		t.Error()
	}
}

func TestLexWhiteSpaceTrailingNumber(t *testing.T) {
	lexer := new("   1")
	lexWhiteSpace(&lexer)
	if lexer.currentIndex != 3 {
		t.Error()
	}
}


func TestLexWhiteSpaceFails(t *testing.T) {
	lexer := new("123 ")
	err := lexWhiteSpace(&lexer)
	if err == nil {
		t.Error()
	}
}

func TestLexWhiteSpaceFailsOnEmpty(t *testing.T) {
	lexer := new("")
	err := lexWhiteSpace(&lexer)
	if err == nil {
		t.Error()
	}
}

func TestLexCase1(t *testing.T) {
	result, _ := Lex("123 123 ( *")
	expected := expectedLexItems(
		[]Token{Number, Number, Parentheses, Operator},
		[]string{"123", "123", "(", "*"},
	)
	if !reflect.DeepEqual(result, expected) {
		t.Error()
	}
}

func TestLexCase2(t *testing.T) {
	result, _ := Lex("  \n\t 1234*12+(123\t*2)  \n")
	expected := expectedLexItems(
		[]Token{Number, Operator, Number, Operator, Parentheses, Number, Operator, Number, Parentheses},
		[]string{"1234", "*", "12", "+", "(", "123", "*", "2", ")"},
	)
	if !reflect.DeepEqual(result, expected) {
		t.Error()
	}
}

func TestLexCase3(t *testing.T) {
	result, _ := Lex("1+(2*3)+4")
	expected := expectedLexItems(
		[]Token{Number, Operator, Parentheses, Number, Operator, Number, Parentheses, Operator, Number},
		[]string{"1", "+", "(", "2", "*", "3", ")", "+", "4"},
	)
	if !reflect.DeepEqual(result, expected) {
		t.Error()
	}
}
