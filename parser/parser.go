package parser

import (
	"fmt"
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

func (parser *Parser) Parse() (grammar.Exp, error) {
	kind, token := parser.readIgnoreWhiteSpace()

	switch kind {
	case tokenizer.ParenthesesStart:
		parser.unread()
		return parser.ParseExpressionStartingWithParentheses()
	case tokenizer.Number:
		parser.unread()
		return parser.ParseExpressionStartingWithNumber()
	}

	return nil, fmt.Errorf("unexpected token, expected parantheses or number: (%d, %s)", kind, token)
}

func (parser *Parser) ParseExpressionStartingWithParentheses() (grammar.Exp, error) {
	// Read open parentheses
	parser.readIgnoreWhiteSpace()

	// Read expression inside parentheses
	expr, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("error while parsing expression inside parentheses: %w", err)
	}

	// Read closing parentheses
	endKind, _ := parser.readIgnoreWhiteSpace()
	if endKind != tokenizer.ParenthesesEnd {
		return nil, fmt.Errorf("closing parentheses missing")
	}

	parenthesesExpression := grammar.ExpParentheses{Inside: expr}

	// Check if the expression is the first part of a binary operation
	nextKind, _ := parser.readIgnoreWhiteSpace()

	if nextKind == tokenizer.Plus {
		nextExpr, err := parser.Parse()
		return grammar.ExpPlus{Left: parenthesesExpression, Right: nextExpr}, err
	}

	if nextKind == tokenizer.Multiply {
		nextExpr, err := parser.Parse()
		return grammar.ExpMultiply{Left: parenthesesExpression, Right: nextExpr}, err
	}

	parser.unread()

	return parenthesesExpression, err
}

func (parser *Parser) ParseExpressionStartingWithNumber() (grammar.Exp, error) {
	_, firstNumber := parser.readIgnoreWhiteSpace()
	firstValue, _ := strconv.Atoi(firstNumber)
	nextKind, _ := parser.readIgnoreWhiteSpace()
	if nextKind == tokenizer.Plus {
		nextExpression, err := parser.Parse()
		return grammar.ExpPlus{
			Left: grammar.ExpNum{Value: firstValue},
			Right: nextExpression,
		}, err
	} else if nextKind == tokenizer.Multiply {
		nextExpression, err := parser.Parse()
		return grammar.ExpMultiply{
			Left: grammar.ExpNum{Value: firstValue},
			Right: nextExpression,
		}, err
	}
	parser.unread()
	return grammar.ExpNum{Value: firstValue}, nil
}
