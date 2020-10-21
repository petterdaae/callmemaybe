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

func (parser *Parser) Parse() (grammar.Stmt, error) {
	stmt, err := parser.parseSeq()
	nextKind, _ := parser.readIgnoreWhiteSpace()
	if nextKind != tokenizer.EOF {
		return nil, fmt.Errorf("failed to parse the entire program")
	}
	return stmt, err
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
	if nextKind == tokenizer.Let {
		parser.unread()
		return parser.parseLet()
	}
	if nextKind == tokenizer.Identifier {
		return grammar.ExpIdentifier{
			Name: nextToken,
		}, nil
	}
	return nil, fmt.Errorf("unexpected token while parsing val")
}

func (parser *Parser) parseLet() (grammar.Exp, error) {
	letKind, _ := parser.readIgnoreWhiteSpace()
	if letKind != tokenizer.Let {
		return nil, fmt.Errorf("expected let to come first when parsing let expression")
	}

	identKind, identifier := parser.readIgnoreWhiteSpace()
	if identKind != tokenizer.Identifier {
		return nil, fmt.Errorf("expected identifier after let in let expression")
	}

	assign, _ := parser.readIgnoreWhiteSpace()
	if assign != tokenizer.Assign {
		return nil, fmt.Errorf("expexted = after identifier in let expression")
	}

	exprIdent, err := parser.parseExp()
	if err != nil {
		return nil, fmt.Errorf("failed to parse first expression in let expression")
	}

	in, _ := parser.readIgnoreWhiteSpace()
	if in != tokenizer.In {
		return nil, fmt.Errorf("expected keyword in after first expression in let expression")
	}

	expr, err := parser.parseExp()
	if err != nil {
		return nil, fmt.Errorf("failed to parse last expression in let expression")
	}

	let := grammar.ExpLet{
		Identifier: identifier,
		IdentifierExp: exprIdent,
		Inside: expr,
	}

	return let, nil
}

func (parser *Parser) parseAssign() (grammar.Stmt, error) {
	kind, identifier := parser.readIgnoreWhiteSpace()
	if kind != tokenizer.Identifier {
		return nil, fmt.Errorf("failed to parse identifier at start of assign statement")
	}
	kind, token := parser.readIgnoreWhiteSpace()
	if kind != tokenizer.Assign {
		return nil, fmt.Errorf("expected assign operator in assign stmt but got: %s", token)
	}
	expr, err := parser.parseExp()
	if err != nil {
		return nil, fmt.Errorf("failed to parse expression in assign stmt: %w", err)
	}
	return grammar.StmtAssign{Identifier: identifier, Expression: expr}, nil
}

func (parser *Parser) parsePrintln() (grammar.Stmt, error) {
	kind, _ := parser.readIgnoreWhiteSpace()
	if kind != tokenizer.PrintLn {
		return nil, fmt.Errorf("expected println keyword at start of println stmt")
	}
	expr, err := parser.parseExp()
	if err != nil {
		return nil, fmt.Errorf("failed to parse expression in print stmt: %w", err)
	}
	return grammar.StmtPrintln{Expression: expr}, nil
}

func (parser *Parser) parseSeq() (grammar.Stmt, error) {
	var statements []grammar.Stmt
	for {
		nextKind, _ := parser.readIgnoreWhiteSpace()
		if nextKind == tokenizer.Identifier {
			parser.unread()
			statement, err := parser.parseAssign()
			if err != nil {
				return nil, fmt.Errorf("failed to parse assign expression: %w", err)
			}
			statements = append(statements, statement)
			continue
		}
		if nextKind == tokenizer.PrintLn {
			parser.unread()
			statement, err := parser.parsePrintln()
			if err != nil {
				return nil, fmt.Errorf("failed to parse println expression: %w", err)
			}
			statements = append(statements, statement)
			continue
		}
		parser.unread()
		break
	}
	return grammar.StmtSeq{Statements: statements}, nil
}
