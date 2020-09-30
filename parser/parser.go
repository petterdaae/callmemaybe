package parser

import (
	"fmt"
	"lang/grammar"
	"lang/lexer"
	"strconv"
)

type parserState struct {
	program           []lexer.LexItem
	currentIndex      int
	currentExpression grammar.Exp
}

func Parse(program []lexer.LexItem) (grammar.Exp, error) {
	state := parserState{
		program: program,
		currentIndex: 0,
		currentExpression: nil,
	}

	for {
		if state.currentIndex >= len(state.program) {
			return state.currentExpression, nil // empty program or end of program
		}
		exp, _ := next(&state)
		state.currentExpression = exp
	}
	return state.currentExpression, nil
}

func next(state *parserState) (grammar.Exp, error) {
	current := state.program[state.currentIndex]

	if current.Kind == lexer.Number {
		exp, err := parseExpNum(state)
		state.currentIndex++
		return exp, err
	}

	switch current.Value() {
	case "(":
		state.currentIndex++
		exp, _ := next(state)
		return grammar.ExpParentheses{Inside: exp}, nil
	case "+":
		state.currentIndex++
		exp, _ := next(state)
		return grammar.ExpPlus{Left: state.currentExpression, Right: exp}, nil
	case "*":
		state.currentIndex++
		exp, _ := next(state)
		return grammar.ExpMultiply{Left: state.currentExpression, Right: exp}, nil
	case ")":
		state.currentIndex++
	}

	return state.currentExpression, fmt.Errorf("invalid next lex item")
}

func parseExpNum(state *parserState) (grammar.ExpNum, error) {
	item := state.program[state.currentIndex]
	value, err := strconv.Atoi(item.Value())
	return grammar.ExpNum{Value: value}, err
}
