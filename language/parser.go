package language

import (
	"callmemaybe/language/typesystem"
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

func (parser *Parser) parseString() (Exp, error) {
	kind, value := parser.readIgnoreWhiteSpace()
	if kind != String {
		return nil, fmt.Errorf("exptected strings kind when parsing string")
	}

	list := ExpList{
		Type: typesystem.Type{
			RawType: typesystem.List,
			ListElementType: &typesystem.Type{
				RawType: typesystem.Char,
			},
		},
		Size: len(value),
	}

	for i := range value {
		char := ExpChar{
			Value: string(value[i]),
		}
		list.Elements = append(list.Elements, char)
	}

	return list, nil
}

func (parser *Parser) ParseExp() (Exp, error) {
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
		if nextKind == Divide {
			right, err := parser.parseVal()
			if err != nil {
				return nil, fmt.Errorf("failed to parse right side of divide exp: %w", err)
			}
			left = ExpDivide{
				Left:  left,
				Right: right,
			}
			continue
		}
		if nextKind == Modulo {
			right, err := parser.parseVal()
			if err != nil {
				return nil, fmt.Errorf("failed to parse right side of modulo exp: %w", err)
			}
			left = ExpModulo{
				Left:  left,
				Right: right,
			}
			continue
		}
		if nextKind == Minus {
			right, err := parser.parseVal()
			if err != nil {
				return nil, fmt.Errorf("failed to parse right side of minus exp: %w", err)
			}
			left = ExpMinus{
				Left:  left,
				Right: right,
			}
			continue
		}
		if nextKind == AngleBracketStart {
			right, err := parser.parseVal()
			if err != nil {
				return nil, fmt.Errorf("failed to parse right side of less expression")
			}
			left = ExpLess{
				Left:  left,
				Right: right,
			}
			continue
		}
		if nextKind == AngleBracketEnd {
			right, err := parser.parseVal()
			if err != nil {
				return nil, fmt.Errorf("failed to parse right side of greater expression")
			}
			left = ExpGreater{
				Left:  left,
				Right: right,
			}
			continue
		}
		if nextKind == Equals {
			right, err := parser.parseVal()
			if err != nil {
				return nil, fmt.Errorf("failed to parse right side of equals expression")
			}
			left = ExpEquals{
				Left:  left,
				Right: right,
			}
			continue
		}
		if nextKind == NotEqual {
			right, err := parser.parseVal()
			if err != nil {
				return nil, fmt.Errorf("failed to parse right side of not equals expression")
			}
			left = ExpNotEquals{
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
	if nextKind == True {
		return ExpBool{
			Value: true,
		}, nil
	}
	if nextKind == False {
		return ExpBool{
			Value: false,
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
	if nextKind == Minus {
		inside, err := parser.ParseExp()
		if err != nil {
			return nil, fmt.Errorf("failed to parse exp in negative expression: %w", err)
		}
		return ExpNegative{
			Inside: inside,
		}, nil
	}
	if nextKind == Character {
		return ExpChar{
			Value: nextToken,
		}, nil
	}
	if nextKind == Question {
		parser.unread()
		return parser.parseGetFromList()
	}
	if nextKind == Hash {
		parser.unread()
		return parser.parseCall()
	}
	if nextKind == Length {
		parser.unread()
		return parser.parseLength()
	}
	if nextKind == Pipe {
		parser.unread()
		return parser.parseFunction()
	}
	if nextKind == AngleBracketStart {
		parser.unread()
		return parser.parseList()
	}
	if nextKind == String {
		parser.unread()
		return parser.parseString()
	}
	if nextKind == At {
		parser.unread()
		return parser.parseStructInit()
	}
	return nil, fmt.Errorf("unexpected token while parsing val")
}

func (parser *Parser) parseLength() (Exp, error) {
	kind, _ := parser.readIgnoreWhiteSpace()
	if kind != Length {
		return nil, fmt.Errorf("expected length keyword")
	}
	kind, _ = parser.readIgnoreWhiteSpace()
	if kind != RoundBracketStart {
		return nil, fmt.Errorf("expected (")
	}
	exp, err := parser.ParseExp()
	if err != nil {
		return nil, fmt.Errorf("exp in lenght: %w", err)
	}

	kind, _ = parser.readIgnoreWhiteSpace()
	if kind != RoundBracketEnd {
		return nil, fmt.Errorf("expected )")
	}
	return ExpLength{List: exp}, nil
}

func (parser *Parser) parseAssign() (Stmt, error) {
	kind, identifier := parser.readIgnoreWhiteSpace()
	if kind != Identifier && kind != Placeholder {
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
		if nextKind == Identifier || nextKind == Placeholder {
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
		if nextKind == If {
			parser.unread()
			statement, err := parser.parseIf()
			if err != nil {
				return nil, fmt.Errorf("failed to parse if statement: %w", err)
			}
			statements = append(statements, statement)
			continue
		}
		if nextKind == Loop {
			parser.unread()
			statement, err := parser.parseLoop()
			if err != nil {
				return nil, fmt.Errorf("failed to parse loop: %w", err)
			}
			statements = append(statements, statement)
			continue
		}
		if nextKind == Struct {
			parser.unread()
			statement, err := parser.parseStructDeclaration()
			if err != nil {
				return nil, fmt.Errorf("failed to parse struct declaration: %w", err)
			}
			statements = append(statements, statement)
			continue
		}
		if nextKind == Question {
			parser.unread()
			statement, err := parser.parseUpdateStmt()
			if err != nil {
				return nil, fmt.Errorf("failed to parse update: %w", err)
			}
			statements = append(statements, statement)
			continue
		}
		parser.unread()
		break
	}
	return StmtSeq{Statements: statements}, nil
}

func (parser *Parser) parseStructInit() (Exp, error) {
	kind, _ := parser.readIgnoreWhiteSpace()
	if kind != At {
		return nil, fmt.Errorf("expected @")
	}
	kind, name := parser.readIgnoreWhiteSpace()
	if kind != Identifier {
		return nil, fmt.Errorf("expected identifier")
	}
	kind, _ = parser.readIgnoreWhiteSpace()
	if kind != CurlyBracketStart {
		return nil, fmt.Errorf("expected {")
	}
	structExp := StructExp{
		Name:    name,
	}
	for {
		kind, memberName := parser.readIgnoreWhiteSpace()
		if kind == CurlyBracketEnd {
			break
		}
		if kind != Identifier {
			return nil, fmt.Errorf("expected identifier")
		}
		kind, _ = parser.readIgnoreWhiteSpace()
		if kind != Colon {
			return nil, fmt.Errorf("expected colon")
		}
		exp, err := parser.ParseExp()
		if err != nil {
			return nil, fmt.Errorf("failed to parse expression in struct member")
		}
		structExp.Members = append(structExp.Members, StructMember{
			Name: memberName,
			Exp:  exp,
		})
	}
	return structExp, nil
}

func (parser *Parser) parseLoop() (Stmt, error) {
	kind, _ := parser.readIgnoreWhiteSpace()
	if kind != Loop {
		return nil, fmt.Errorf("expected loop keyword")
	}

	exp, err := parser.ParseExp()
	if err != nil {
		return nil, fmt.Errorf("loop condition: %w", err)
	}

	kind, _ = parser.readIgnoreWhiteSpace()
	if kind != CurlyBracketStart {
		return nil, fmt.Errorf("expected { in loop")
	}

	body, err := parser.parseSeq()
	if err != nil {
		return nil, fmt.Errorf("loop body: %w", err)
	}

	kind, _ = parser.readIgnoreWhiteSpace()
	if kind != CurlyBracketEnd {
		return nil, fmt.Errorf("expected } in loop")
	}

	return StmtLoop{
		Condition: exp,
		Body:      body,
	}, nil
}

func (parser *Parser) parseIf() (Stmt, error) {
	kind, _ := parser.readIgnoreWhiteSpace()
	if kind != If {
		return nil, fmt.Errorf("expected if keyword at start of if statement")
	}

	expr, err := parser.ParseExp()
	if err != nil {
		return nil, fmt.Errorf("failed to parse condition expression in if statement: %w", err)
	}

	kind, text := parser.readIgnoreWhiteSpace()
	if kind != CurlyBracketStart {
		return nil, fmt.Errorf("expected { when parsing if statement, but got: %s", text)
	}

	seq, err := parser.parseSeq()
	if err != nil {
		return nil, fmt.Errorf("failed to parse sequence in if statement: %w", err)
	}

	kind, text = parser.readIgnoreWhiteSpace()
	if kind != CurlyBracketEnd {
		return nil, fmt.Errorf("expected } when parsing if statement, but got: %s", text)
	}

	return StmtIf{
		Expression: expr,
		Body:       seq,
	}, nil
}

func (parser *Parser) parseCall() (Exp, error) {
	kind, _ := parser.readIgnoreWhiteSpace()
	if kind != Hash {
		return nil, fmt.Errorf("expected # in call")
	}

	expr, err := parser.ParseExp()
	if err != nil {
		return nil, fmt.Errorf("expression in call: %w", err)
	}

	call := FunctionCall{}
	call.Exp = expr

	kind, _ = parser.readIgnoreWhiteSpace()
	if kind == RoundBracketStart {
		for {
			expr, err := parser.ParseExp()
			if err != nil {
				return nil, fmt.Errorf("failed to parse expression in function call")
			}
			call.Arguments = append(call.Arguments, expr)
			kind, _ = parser.readIgnoreWhiteSpace()
			if kind != Comma {
				parser.unread()
				break
			}

		}
		kind, _ = parser.readIgnoreWhiteSpace()
		if kind != RoundBracketEnd {
			return nil, fmt.Errorf("expected ) in call")
		}
	} else {
		parser.unread()
	}

	return call, nil
}

func (parser *Parser) parseFunction() (Exp, error) {
	function := ExpFunction{
		Type: typesystem.Type{
			RawType: typesystem.Function,
		},
	}
	kind, _ := parser.readIgnoreWhiteSpace()
	if kind != Pipe {
		return nil, fmt.Errorf("expected |")
	}
	first := true
	for {
		kind, identifier := parser.readIgnoreWhiteSpace()
		if kind == Pipe {
			break
		}
		if kind != Identifier {
			return nil, fmt.Errorf("expected identifier")
		}
		kind, _ = parser.readIgnoreWhiteSpace()
		if kind == Comma && first {
			function.Recurse = identifier
			if identifier != "me" {
				return nil, fmt.Errorf("the first recurse argument in a function has to be named me")
			}
			continue
		} else {
			parser.unread()
		}
		first = false

		argType, err := parser.parseType()
		if err != nil {
			return nil, fmt.Errorf("failed to parse type: %w", err)
		}
		if !argType.IsPassable() {
			return nil, fmt.Errorf("expected passable type when parsing function arguments")
		}
		function.Type.FunctionArgumentTypes = append(function.Type.FunctionArgumentTypes, typesystem.NamedType{
			Name: identifier,
			Type: argType,
		})
		kind, _ = parser.readIgnoreWhiteSpace()
		if kind == Pipe {
			break
		}
		if kind == Comma {
			continue
		}
		return nil, fmt.Errorf("expected comma or end of argument list")
	}

	kind, _ = parser.readIgnoreWhiteSpace()
	if kind == CurlyBracketStart {
		function.Type.FunctionReturnType = &typesystem.Type{
			RawType: typesystem.Void,
		}
	} else {
		parser.unread()
		returnType, err := parser.parseType()
		if err != nil {
			return nil, fmt.Errorf("failed to parse function return type")
		}
		if !returnType.IsPassable() {
			return nil, fmt.Errorf("expected passable type when parsing function return type")
		}
		function.Type.FunctionReturnType = &returnType
		kind, _ = parser.readIgnoreWhiteSpace()
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

func (parser *Parser) parseList() (Exp, error) {
	kind, _ := parser.readIgnoreWhiteSpace()
	if kind != AngleBracketStart {
		return nil, fmt.Errorf("expected angle bracket when parsing list")
	}
	_type, err := parser.parseType()
	if err != nil {
		return nil, fmt.Errorf("failed to parse list type: %w", err)
	}
	kind, _ = parser.readIgnoreWhiteSpace()
	if kind != Comma {
		return nil, fmt.Errorf("expected comma when parsing list type")
	}
	kind, numberStr := parser.readIgnoreWhiteSpace()
	if kind != Number {
		return nil, fmt.Errorf("list size should be a number")
	}
	number, _ := strconv.Atoi(numberStr)
	kind, _ = parser.readIgnoreWhiteSpace()
	if kind != AngleBracketEnd {
		return nil, fmt.Errorf("expected closing angle bracket")
	}
	kind, _ = parser.readIgnoreWhiteSpace()
	if kind != BoxBracketStart {
		return nil, fmt.Errorf("expected box bracket in list declaration")
	}

	list := ExpList{}
	first := true

	for {
		kind, _ = parser.readIgnoreWhiteSpace()
		if first && kind == BoxBracketEnd {
			break
		}
		first = false
		parser.unread()
		expr, err := parser.ParseExp()
		if err != nil {
			return nil, fmt.Errorf("failed to parse expression in list declaration: %w", err)
		}
		list.Elements = append(list.Elements, expr)
		kind, _ = parser.readIgnoreWhiteSpace()
		if kind == Comma {
			continue
		}
		if kind == BoxBracketEnd {
			break
		}
		return nil, fmt.Errorf("unexpected token when parsing list declareation")
	}

	list.Type = typesystem.Type{
		RawType:         typesystem.List,
		ListElementType: &_type,
	}

	list.Size = number

	return list, nil
}

func (parser *Parser) parseGetFromList() (Exp, error) {
	return parser.parseReference()
}

func (parser *Parser) parseStructDeclaration() (Stmt, error) {
	kind, _ := parser.readIgnoreWhiteSpace()
	if kind != Struct {
		return nil, fmt.Errorf("expected struct keyword")
	}
	kind, name := parser.readIgnoreWhiteSpace()
	if kind != Identifier {
		return nil, fmt.Errorf("expected identifier")
	}
	kind, _ = parser.readIgnoreWhiteSpace()
	if kind != CurlyBracketStart {
		return nil, fmt.Errorf("expected curly bracket")
	}
	structType := typesystem.Type{
		RawType:               typesystem.Struct,
		StructName:            name,
	}
	for {
		kind, memberName := parser.readIgnoreWhiteSpace()
		if kind == CurlyBracketEnd {
			break
		}
		if kind != Identifier {
			return nil, fmt.Errorf("expected identifier")
		}
		_type, err := parser.parseType()
		if err != nil {
			return nil, fmt.Errorf("struct type: %w", err)
		}
		structType.StructMembers = append(structType.StructMembers, typesystem.NamedType{
			Name: memberName,
			Type: _type,
		})
	}
	return StmtStructDeclaration{
		Type: structType,
	}, nil
}

func (parser *Parser) parseUpdateStmt() (Stmt, error) {
	reference, err := parser.parseReference()
	if err != nil {
		return nil, fmt.Errorf("reference: %w", err)
	}

	kind, _ := parser.readIgnoreWhiteSpace()
	if kind != Assign {
		return nil, fmt.Errorf("expected =")
	}

	newValue, err := parser.ParseExp()
	if err != nil {
		return nil, fmt.Errorf("new value: %w", err)
	}

	if val, ok := reference.(ExpReadFromStruct); ok {
		return StmtUpdateStruct{
			Struct:   val.Struct,
			Member:   val.Field,
			NewValue: newValue,
		}, nil
	}

	if val, ok := reference.(ExpGetFromList); ok {
		return StmtUpdateList{
			List:     val.List,
			Index:    val.Index,
			NewValue: newValue,
		}, nil
	}

	return nil, fmt.Errorf("invalid reference")
}

func (parser *Parser) parseReference() (Exp, error) {
	kind, _ := parser.readIgnoreWhiteSpace()
	if kind != Question {
		return nil, fmt.Errorf("expected ? ")
	}
	exp, err := parser.ParseExp()
	if err != nil {
		return nil, fmt.Errorf("failed to parse expression in get from list: %w", err)
	}

	current := exp

	hasOne := false

	for {
		kind, _ = parser.readIgnoreWhiteSpace()
		if kind == BoxBracketStart {
			hasOne = true
			numExp, err := parser.ParseExp()
			if err != nil {
				return nil, fmt.Errorf("failed to parse index expression in get from list")
			}
			kind, _ = parser.readIgnoreWhiteSpace()
			if kind != BoxBracketEnd {
				return nil, fmt.Errorf("expected ]")
			}
			current = ExpGetFromList{
				List:  current,
				Index: numExp,
			}
			continue
		}
		if kind == Dot {
			hasOne = true
			kind, identifier := parser.readIgnoreWhiteSpace()
			if kind != Identifier {
				return nil, fmt.Errorf("expected identifier")
			}
			current = ExpReadFromStruct{
				Field:  identifier,
				Struct: current,
			}
			continue
		}
		parser.unread()
		break
	}

	if !hasOne {
		return nil, fmt.Errorf("expected . or [")
	}

	return current, nil
}

func (parser *Parser) parseType() (typesystem.Type, error) {
	kind, _ := parser.readIgnoreWhiteSpace()
	switch kind {
	case TypeInt:
		return typesystem.Type{
			RawType: typesystem.Int,
		}, nil
	case TypeBool:
		return typesystem.Type{
			RawType: typesystem.Bool,
		}, nil
	case TypeChar:
		return typesystem.Type{
			RawType: typesystem.Char,
		}, nil
	case TypeList:
		kind, _ = parser.readIgnoreWhiteSpace()
		if kind != AngleBracketStart {
			return typesystem.Type{}, fmt.Errorf("expected opening angle bracket")
		}
		elementType, err := parser.parseType()
		if err != nil {
			return typesystem.Type{}, fmt.Errorf("failed to parse list element type: %w", err)
		}
		kind, _ = parser.readIgnoreWhiteSpace()
		if kind != AngleBracketEnd {
			return typesystem.Type{}, fmt.Errorf("expected closing angle bracket when parsing list type")
		}
		return typesystem.Type{
			RawType:               typesystem.List,
			ListElementType:       &elementType,
		}, nil
	case TypeString:
		return typesystem.Type{
			RawType:               typesystem.List,
			ListElementType:       &typesystem.Type{
				RawType:               typesystem.Char,
			},
		}, nil
	case TypeFunc:
		kind, _ = parser.readIgnoreWhiteSpace()
		if kind != AngleBracketStart {
			parser.unread()
			return typesystem.Type{
				RawType:               typesystem.Function,
				FunctionArgumentTypes: nil,
				FunctionReturnType:    &typesystem.Type{
					RawType:               typesystem.Void,
				},
			}, nil
		}
		var types []typesystem.Type
		for {
			_type, err := parser.parseType()
			if err != nil {
				return typesystem.Type{}, fmt.Errorf("type inside function type: %w", err)
			}
			types = append(types, _type)
			kind, _ = parser.readIgnoreWhiteSpace()
			if kind == Comma {
				continue
			}
			parser.unread()
			break
		}
		kind, _ = parser.readIgnoreWhiteSpace()
		if kind != AngleBracketEnd {
			return typesystem.Type{}, fmt.Errorf("expected > at end of function tye")
		}
		size := len(types)
		result := typesystem.Type{
			RawType:               typesystem.Function,
		}
		for i := 0; i < size - 1; i++ {
			current := types[i]
			result.FunctionArgumentTypes = append(result.FunctionArgumentTypes, typesystem.NamedType{
				Name: "",
				Type: current,
			})
		}
		result.FunctionReturnType = &types[size-1]
		return result, nil
	case At:
		kind, name := parser.readIgnoreWhiteSpace()
		if kind != Identifier {
			return typesystem.NewInvalid(), fmt.Errorf("expected identifier")
		}
		return typesystem.Type{
			RawType:               typesystem.Struct,
			StructName:            name,
		}, nil
	default:
		return typesystem.Type{}, fmt.Errorf("unsupported type")
	}
}
