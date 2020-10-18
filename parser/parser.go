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

func (parser *Parser) Parse() (grammar.Exp, error) {
	expr, err := parser.parseExp()
	nextKind, _ := parser.readIgnoreWhiteSpace()
	if nextKind != tokenizer.EOF {
		return nil, fmt.Errorf("failed to parse the entire program")
	}
	return expr, err
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

func (parser *Parser) parseExp() (grammar.Exp, error) {
	left, err := parser.parseVal()
	if err != nil {
		return nil, fmt.Errorf("failed to parse first val in exp: %w", err)
	}
	for {
		nextKind, _ := parser.readIgnoreWhiteSpace()
		if nextKind == tokenizer.Plus {
			right, err := parser.parseVal()
			if err != nil {
				return nil, fmt.Errorf("failed to parse right side of plus exp: %w", err)
			}
			left = grammar.ExpPlus{
				Left:  left,
				Right: right,
			}
			continue
		}
		if nextKind == tokenizer.Multiply {
			right, err := parser.parseVal()
			if err != nil {
				return nil, fmt.Errorf("failed to parse right side of multiply exp: %w", err)
			}
			left = grammar.ExpMultiply{
				Left:  left,
				Right: right,
			}
			continue
		}
		parser.unread()
		break
	}
	return left, nil
}

func (parser *Parser) parseVal() (grammar.Exp, error) {
	nextKind, nextToken := parser.readIgnoreWhiteSpace()
	if nextKind == tokenizer.Number {
		value, _ := strconv.Atoi(nextToken)
		return grammar.ExpNum{
			Value: value,
		}, nil
	}
	if nextKind == tokenizer.ParenthesesStart {
		inside, err := parser.parseExp()
		if err != nil {
			return nil, fmt.Errorf("failed to parse exp in parentheses: %w", err)
		}
		nextKind, _ = parser.readIgnoreWhiteSpace()
		if nextKind != tokenizer.ParenthesesEnd {
			return nil, fmt.Errorf("missing closing parentheses")
		}
		return grammar.ExpParentheses{
			Inside: inside,
		}, nil
	}
	return nil, fmt.Errorf("unexpected token while parsing val")
}
