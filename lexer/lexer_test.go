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
