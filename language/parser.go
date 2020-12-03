package language

import (
	"fmt"
	"io"
	"strconv"
)

type Parser struct {
	tokenizer *Tokenizer
	buffer    struct {
		kind  Token
		token string
		full  bool
	}
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{
		tokenizer: NewTokenizer(reader),
	}
}

func (parser *Parser) Parse() (Stmt, error) {
	stmt, err := parser.parseSeq()
	nextKind, nextStr := parser.readIgnoreWhiteSpace()
	if err == nil && nextKind != EOF {
		return nil, fmt.Errorf("failed to parse the entire program: %s", nextStr)
	}
	return stmt, err
}

func (parser *Parser) read() (Token, string) {
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

func (parser *Parser) readIgnoreWhiteSpace() (Token, string) {
	kind, token := parser.read()
	if kind == Whitespace {
		kind, token = parser.read()
	}
	return kind, token
}

func (parser *Parser) ParseExp() (Exp, error) {
	nextKind, _ := parser.readIgnoreWhiteSpace()
	parser.unread()

	if nextKind == Call {
		return parser.parseCall()
	}

	if nextKind == AngleBracketStart {
		return parser.parseFunction()
	}

	return parser.ParseCalculation()
}

func (parser *Parser) ParseCalculation() (Exp, error) {
	left, err := parser.parseVal()
	if err != nil {
		return nil, fmt.Errorf("failed to parse first val in exp: %w", err)
	}
	for {
		nextKind, _ := parser.readIgnoreWhiteSpace()
		if nextKind == Plus {
			right, err := parser.parseVal()
			if err != nil {
				return nil, fmt.Errorf("failed to parse right side of plus exp: %w", err)
			}
			left = ExpPlus{
				Left:  left,
				Right: right,
			}
			continue
		}
		if nextKind == Multiply {
			right, err := parser.parseVal()
			if err != nil {
				return nil, fmt.Errorf("failed to parse right side of multiply exp: %w", err)
			}
			left = ExpMultiply{
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

func (parser *Parser) parseVal() (Exp, error) {
	nextKind, nextToken := parser.readIgnoreWhiteSpace()
	if nextKind == Number {
		value, _ := strconv.Atoi(nextToken)
		return ExpNum{
			Value: value,
		}, nil
	}
	if nextKind == RoundBracketStart {
		inside, err := parser.ParseExp()
		if err != nil {
			return nil, fmt.Errorf("failed to parse exp in parentheses: %w", err)
		}
		nextKind, _ = parser.readIgnoreWhiteSpace()
		if nextKind != RoundBracketEnd {
			return nil, fmt.Errorf("missing closing parentheses")
		}
		return ExpParentheses{
			Inside: inside,
		}, nil
	}
	if nextKind == Identifier {
		return ExpIdentifier{
			Name: nextToken,
		}, nil
	}
	return nil, fmt.Errorf("unexpected token while parsing val")
}

func (parser *Parser) parseAssign() (Stmt, error) {
	kind, identifier := parser.readIgnoreWhiteSpace()
	if kind != Identifier {
		return nil, fmt.Errorf("failed to parse identifier at start of assign statement")
	}
	kind, token := parser.readIgnoreWhiteSpace()
	if kind != Assign {
		return nil, fmt.Errorf("expected assign operator in assign stmt but got: %s", token)
	}
	expr, err := parser.ParseExp()
	if err != nil {
		return nil, fmt.Errorf("failed to parse expression in assign stmt: %w", err)
	}
	return StmtAssign{Identifier: identifier, Expression: expr}, nil
}

func (parser *Parser) parsePrintln() (Stmt, error) {
	kind, _ := parser.readIgnoreWhiteSpace()
	if kind != PrintLn {
		return nil, fmt.Errorf("expected println keyword at start of println stmt")
	}
	expr, err := parser.ParseExp()
	if err != nil {
		return nil, fmt.Errorf("failed to parse expression in print stmt: %w", err)
	}
	return StmtPrintln{Expression: expr}, nil
}

func (parser *Parser) parseSeq() (Stmt, error) {
	var statements []Stmt
	for {
		nextKind, _ := parser.readIgnoreWhiteSpace()
		if nextKind == Identifier {
			parser.unread()
			statement, err := parser.parseAssign()
			if err != nil {
				return nil, fmt.Errorf("failed to parse assign expression: %w", err)
			}
			statements = append(statements, statement)
			continue
		}
		if nextKind == PrintLn {
			parser.unread()
			statement, err := parser.parsePrintln()
			if err != nil {
				return nil, fmt.Errorf("failed to parse println expression: %w", err)
			}
			statements = append(statements, statement)
			continue
		}
		if nextKind == Return {
			expr, err := parser.ParseExp()
			if err != nil {
				return nil, fmt.Errorf("failed to parse expression after return: %w", err)
			}
			statement := StmtReturn{Expression: expr}
			statements = append(statements, statement)
			continue
		}
		parser.unread()
		break
	}
	return StmtSeq{Statements: statements}, nil
}

func (parser *Parser) parseCall() (ExprStmt, error) {
	// TODO : parse call
	return StmtSeq{}, nil
}

func (parser *Parser) parseFunction() (Exp, error) {
	function := ExpFunction{}
	kind, _ := parser.readIgnoreWhiteSpace()
	if kind != AngleBracketStart {
		return nil, fmt.Errorf("expected < at start of function expression")
	}

	for {
		kind, identifier := parser.readIgnoreWhiteSpace()
		if kind == AngleBracketEnd {
			break
		}

		if kind != Identifier {
			return nil, fmt.Errorf("expected identifier when parsing argument list, but got %s", identifier)
		}

		kind, typeName := parser.readIgnoreWhiteSpace()
		if kind == TypeEmpty {
			return nil, fmt.Errorf("only non-empty types allowed in function arguments")
		}

		if kind != TypeInt {
			return nil, fmt.Errorf("expected non-empty type as argument type when parsing argument list")
		}

		kind, _ = parser.readIgnoreWhiteSpace()

		function.Args = append(function.Args, Arg{Identifier: identifier, Type: typeName})

		if kind == AngleBracketEnd {
			break
		}

		if kind == Comma {
			continue
		}

		return nil, fmt.Errorf("expected comma or end of arguemnt list")
	}

	kind, _ = parser.readIgnoreWhiteSpace()
	if kind != Arrow {
		return nil, fmt.Errorf("expected arrow after argument list when parsing function")
	}

	kind, typeName := parser.readIgnoreWhiteSpace()
	if kind == TypeInt {
		function.ReturnType = typeName
		kind, _ = parser.readIgnoreWhiteSpace()
	} else {
		function.ReturnType = "empty"
	}

	if kind != CurlyBracketStart {
		return nil, fmt.Errorf("expected opening curly bracket when parsing function")
	}

	seq, err := parser.parseSeq()
	if err != nil {
		return nil, fmt.Errorf("failed to parse statements in function: %w", err)
	}

	function.Body = seq

	kind, _ = parser.readIgnoreWhiteSpace()
	if kind != CurlyBracketEnd {
		return nil, fmt.Errorf("expected closing curly bracker when parsing function")
	}

	return function, nil
}
