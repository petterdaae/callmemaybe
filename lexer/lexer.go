package lexer

import (
	"fmt"
	"unicode"
)

type LexItemKind int

type Lexer struct {
	program      []rune
	currentIndex int
	lexItems     []LexItem
	operators    []rune
}

type LexItem struct {
	kind  LexItemKind
	value []rune
}

type LexFunc = func(*Lexer) error

const (
	Number LexItemKind = iota
	Operator
	Parentheses
)

func New(program string) Lexer {
	return Lexer{
		program:      []rune(program),
		currentIndex: 0,
		lexItems:     []LexItem{},
		operators:    []rune{'+', '*'},
	}
}

func Lex(program string) ([]LexItem, error) {
	lexer := New(program)
	for {
		lexers := []func(*Lexer) error{LexNumber, LexOperator, LexParentheses, LexWhiteSpace}
		tryOne := OneOf(lexers)
		err := tryOne(&lexer)
		if err != nil {
			if lexer.currentIndex < len(program) {
				return nil, fmt.Errorf("lexer failed at index %d: %w", lexer.currentIndex, err)
			}
			break
		}
	}
	return lexer.lexItems, nil
}

func OneOf(lexers []LexFunc) LexFunc {
	return func(lexer *Lexer) error {
		for _, lex := range lexers {
			err := lex(lexer)
			if err != nil {
				return nil
			}
		}
		return fmt.Errorf("all lexers failed")
	}
}

func LexOneOfCharacters(chars []rune, kind LexItemKind) LexFunc {
	return func(lexer *Lexer) error {
		if lexer.currentIndex >= len(lexer.program) {
			return fmt.Errorf("can't lex at end of input")
		}
		for _, char := range chars {
			if lexer.program[lexer.currentIndex] == char {
				item := LexItem{
					kind: kind,
					value: []rune{char},
				}
				lexer.lexItems = append(lexer.lexItems, item)
				lexer.currentIndex++
				return nil
			}
		}
		return fmt.Errorf("did not find a character to lex at current index")
	}
}

func LexNumber(lexer *Lexer) error {
	var result []rune
	for i := lexer.currentIndex; i < len(lexer.program); i++ {
		rune := lexer.program[i]
		if unicode.IsDigit(rune) {
			result = append(result, rune)
		} else {
			break
		}
	}
	if len(result) == 0 {
		return fmt.Errorf("did not find a number to lex")
	}
	item := LexItem{
		kind:  Number,
		value: result,
	}
	lexer.lexItems = append(lexer.lexItems, item)
	lexer.currentIndex += len(result)
	return nil
}

func LexOperator(lexer *Lexer) error {
	lex := LexOneOfCharacters(lexer.operators, Operator)
	return lex(lexer)
}

func LexParentheses(lexer *Lexer) error {
	parentheses := []rune{'(', ')'}
	lex := LexOneOfCharacters(parentheses, Parentheses)
	return lex(lexer)
}

func LexWhiteSpace(lexer *Lexer) error {
	return nil
}
