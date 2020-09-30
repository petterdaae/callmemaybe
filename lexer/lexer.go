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
		err := TryOne(&lexer)
		if err != nil {
			if lexer.currentIndex < len(program) {
				return nil, fmt.Errorf("lexer failed at index %d: %w", lexer.currentIndex, err)
			}
			break
		}
	}
	return lexer.lexItems, nil
}

func TryOne(lexer *Lexer) error {
	all := []func(*Lexer) error{LexNumber, LexOperator, LexParentheses, LexWhiteSpace}
	for _, lex := range all {
		err := lex(lexer)
		if err != nil {
			return nil
		}
	}
	return fmt.Errorf("all lexers failed")
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
	if lexer.currentIndex >= len(lexer.program) {
		return fmt.Errorf("can't lex at end of input")
	}

	for _, operator := range lexer.operators {
		if lexer.program[lexer.currentIndex] == operator {
			item := LexItem{
				kind: Operator,
				value: []rune{operator},
			}
			lexer.lexItems = append(lexer.lexItems, item)
			lexer.currentIndex++
			return nil
		}
	}

	return fmt.Errorf("did not find an operator to lex at current index")
}

func LexParentheses(lexer *Lexer) error {
	return nil
}

func LexWhiteSpace(lexer *Lexer) error {
	return nil
}
