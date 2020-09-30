package lexer

import (
	"fmt"
	"unicode"
)

type LexItemKind int

type lexerState struct {
	program      []rune
	currentIndex int
	lexItems     []LexItem
	operators    []rune
	whitespace   []rune
}

type LexItem struct {
	Kind  LexItemKind
	value []rune
}

type lexFunc = func(*lexerState) error

const (
	Number LexItemKind = iota
	Operator
	Parentheses
	Ignore // Whitespace, comments, etc...
)

func (item LexItem) Equals(value string) bool {
	return string(item.value) == value
}

func (item LexItem) Value() string {
	return string(item.value)
}

func (lexer *lexerState) addLexItem(kind LexItemKind, value []rune) {
	if kind != Ignore {
		item := LexItem{
			Kind:  kind,
			value: value,
		}
		lexer.lexItems = append(lexer.lexItems, item)
	}
}

func new(program string) lexerState {
	return lexerState{
		program:      []rune(program),
		currentIndex: 0,
		lexItems:     []LexItem{},
		operators:    []rune{'+', '*'},
		whitespace:   []rune{' ', '\n', '\t'},
	}
}

func Lex(program string) ([]LexItem, error) {
	lexer := new(program)
	lexers := []lexFunc{lexNumber, lexOperator, lexParentheses, lexWhiteSpace}
	tryOne := oneOf(lexers)
	lex := oneOrMany(tryOne)
	err := lex(&lexer)
	if err != nil {
		return nil, fmt.Errorf("lexer failed: %w", err)
	}
	return lexer.lexItems, nil
}

func oneOf(lexers []lexFunc) lexFunc {
	return func(lexer *lexerState) error {
		for _, lex := range lexers {
			err := lex(lexer)
			if err == nil {
				return nil
			}
		}
		return fmt.Errorf("all lexers failed")
	}
}

func oneOrMany(lex lexFunc) lexFunc {
	return func(lexer *lexerState) error {
		count := 0
		for {
			err := lex(lexer)
			if err != nil {
				break
			}
			count++
		}
		if count == 0 {
			return fmt.Errorf("lexer did not run one or many times")
		}
		return nil
	}
}

func lexOneOfCharacters(chars []rune, kind LexItemKind) lexFunc {
	return func(lexer *lexerState) error {
		if lexer.currentIndex >= len(lexer.program) {
			return fmt.Errorf("can't lex at end of input")
		}
		for _, char := range chars {
			lex := lexOneCharacter(char, kind)
			err := lex(lexer)
			if err == nil {
				return nil
			}
		}
		return fmt.Errorf("did not find a character to lex at current index")
	}
}

func lexOneCharacter(char rune, kind LexItemKind) lexFunc {
	return func(lexer *lexerState) error {
		if lexer.currentIndex >= len(lexer.program) {
			return fmt.Errorf("can't lex at end of input")
		}
		if lexer.program[lexer.currentIndex] == char {
			lexer.addLexItem(kind, []rune{char})
			lexer.currentIndex++
			return nil
		}
		return fmt.Errorf("failed to lex character")
	}
}

func lexNumber(lexer *lexerState) error {
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
	lexer.addLexItem(Number, result)
	lexer.currentIndex += len(result)
	return nil
}

func lexOperator(lexer *lexerState) error {
	lex := lexOneOfCharacters(lexer.operators, Operator)
	return lex(lexer)
}

func lexParentheses(lexer *lexerState) error {
	parentheses := []rune{'(', ')'}
	lex := lexOneOfCharacters(parentheses, Parentheses)
	return lex(lexer)
}

func lexWhiteSpace(lexer *lexerState) error {
	lexOne := lexOneOfCharacters(lexer.whitespace, Ignore)
	lexAll := oneOrMany(lexOne)
	return lexAll(lexer)
}
