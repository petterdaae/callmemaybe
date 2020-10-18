package parser

import (
	"io"
	"lang/grammar"
	"lang/tokenizer"
	"strconv"
)

type Parser struct {
	tokenizer *tokenizer.Tokenizer
	buffer    struct {
		kind  tokenizer.Token
		token string
		full  bool
	}
}

func New(reader io.Reader) *Parser {
	return &Parser{
		tokenizer: tokenizer.New(reader),
	}
}

func (parser *Parser) read() (tokenizer.Token, string) {
	if parser.buffer.full {
		parser.buffer.full = false
		return parser.buffer.kind, parser.buffer.token
	}

	kind, token := parser.tokenizer.NextToken()
	parser.buffer.kind = kind
	parser.buffer.token = token

	return kind, token
}

func (parser *Parser) unread() {
	parser.buffer.full = true
}

func (parser *Parser) readIgnoreWhiteSpace() (tokenizer.Token, string) {
	kind, token := parser.read()
	if kind == tokenizer.Whitespace {
		kind, token = parser.read()
	}
	return kind, token
}

func (parser *Parser) Parse() grammar.Exp {
	kind, _ := parser.readIgnoreWhiteSpace()

	switch kind {
	case tokenizer.ParenthesesStart:
		parser.unread()
		return parser.ParseParentheses()
	case tokenizer.Number:
		parser.unread()
		return parser.ParseExpressionStartingWithNumber()
	}

	println("Error")

	return nil
}

func (parser *Parser) ParseParentheses() grammar.Exp {
	parser.readIgnoreWhiteSpace()
	expr := parser.Parse()
	endKind, _ := parser.readIgnoreWhiteSpace()
	if endKind != tokenizer.ParenthesesEnd {
		println("Error: missing closing parentheses")
	}

	temp := grammar.ExpParentheses{Inside: expr}

	nextKind, _ := parser.readIgnoreWhiteSpace()

	if nextKind == tokenizer.Plus || nextKind == tokenizer.Multiply {
		nextExpr := parser.Parse()
		return grammar.ExpPlus{Left: temp, Right: nextExpr}
	}

	return temp
}

func (parser *Parser) ParseExpressionStartingWithNumber() grammar.Exp {
	_, firstNumber := parser.readIgnoreWhiteSpace()
	firstValue, _ := strconv.Atoi(firstNumber)
	nextKind, _ := parser.readIgnoreWhiteSpace()
	if nextKind == tokenizer.Plus {
		nextExpression := parser.Parse()
		return grammar.ExpPlus{
			Left: grammar.ExpNum{Value: firstValue},
			Right: nextExpression,
		}
	} else if nextKind == tokenizer.Multiply {
		nextExpression := parser.Parse()
		return grammar.ExpMultiply{
			Left: grammar.ExpNum{Value: firstValue},
			Right: nextExpression,
		}
	}
	parser.unread()
	return grammar.ExpNum{Value: firstValue}
}
