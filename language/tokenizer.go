package language

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
)

const (
	Number Token = iota
	Character
	Plus
	Minus
	Modulo
	Divide
	Multiply
	Assign
	RoundBracketStart
	RoundBracketEnd
	BoxBracketStart
	BoxBracketEnd
	CurlyBracketStart
	CurlyBracketEnd
	AngleBracketStart
	AngleBracketEnd
	Comma
	Arrow
	Call
	With
	TypeInt
	TypeBool
	TypeChar
	TypeEmpty
	Whitespace
	PrintLn
	Return
	Identifier
	Placeholder
	True
	False
	If
	Equals
	EOF
	Error
)

const eof = rune(0)

type Token int

type Tokenizer struct {
	reader *bufio.Reader
}

func NewTokenizer(reader io.Reader) *Tokenizer {
	return &Tokenizer{
		reader: bufio.NewReader(reader),
	}
}

func (tokenizer *Tokenizer) read() rune {
	character, _, err := tokenizer.reader.ReadRune()
	if err != nil {
		return eof
	}
	return character
}

func (tokenizer *Tokenizer) unread() {
	tokenizer.reader.UnreadRune()
}

func (tokenizer *Tokenizer) NextToken() (Token, string) {
	character := tokenizer.read()

	if unicode.IsSpace(character) {
		tokenizer.unread()
		return tokenizer.whitespace()
	}

	if unicode.IsDigit(character) {
		tokenizer.unread()
		return tokenizer.number()
	}

	if unicode.IsLetter(character) {
		tokenizer.unread()
		return tokenizer.identifier()
	}

	if character == '\'' {
		tokenizer.unread()
		return tokenizer.character()
	}

	switch character {
	case eof:
		return EOF, ""
	case '*':
		return Multiply, string(character)
	case '+':
		return Plus, string(character)
	case '(':
		return RoundBracketStart, string(character)
	case ')':
		return RoundBracketEnd, string(character)
	case '=':
		next := tokenizer.read()
		if next == '>' {
			return Arrow, "=>"
		}
		if next == '=' {
			return Equals, "=="
		}
		tokenizer.unread()
		return Assign, string(character)
	case '{':
		return CurlyBracketStart, string(character)
	case '}':
		return CurlyBracketEnd, string(character)
	case '[':
		return BoxBracketStart, string(character)
	case ']':
		return BoxBracketEnd, string(character)
	case '<':
		return AngleBracketStart, string(character)
	case '>':
		return AngleBracketEnd, string(character)
	case ',':
		return Comma, string(character)
	case '_':
		return Placeholder, string(character)
	case '-':
		return Minus, string(character)
	case '/':
		return Divide, string(character)
	case '%':
		return Modulo, string(character)
	}
	
	return Error, ""
}

func (tokenizer *Tokenizer) character() (Token, string) {
	character := tokenizer.read()
	if character != '\'' {
		return Error, ""
	}

	character = tokenizer.read()

	if character == '\'' {
		return Error, ""
	}

	if character == '\\' {
		character = tokenizer.read()
		if character != '\'' && character != '\\' {
			return Error, ""
		}
	}

	end := tokenizer.read()
	if end != '\'' {
		return Error, ""
	}

	return Character, string(character)
}

func (tokenizer *Tokenizer) whitespace() (Token, string) {
	var buffer bytes.Buffer
	buffer.WriteRune(tokenizer.read())

	for {
		character := tokenizer.read()
		if character == eof || !unicode.IsSpace(character) {
			tokenizer.unread()
			break
		}
		buffer.WriteRune(character)
	}

	return Whitespace, buffer.String()
}

func (tokenizer *Tokenizer) number() (Token, string) {
	var buffer bytes.Buffer
	buffer.WriteRune(tokenizer.read())

	for {
		character := tokenizer.read()
		if !unicode.IsDigit(character) {
			tokenizer.unread()
			break
		}
		buffer.WriteRune(character)
	}
	return Number, buffer.String()
}

func (tokenizer *Tokenizer) identifier() (Token, string) {
	var buffer bytes.Buffer
	buffer.WriteRune(tokenizer.read())

	for {
		character := tokenizer.read()
		if !unicode.IsLetter(character) && !unicode.IsDigit(character){
			tokenizer.unread()
			break
		}
		buffer.WriteRune(character)
	}

	word := buffer.String()

	switch word {
	case "println":
		return PrintLn, word
	case "call":
		return Call, word
	case "with":
		return With, word
	case "int":
		return TypeInt, word
	case "char":
		return TypeChar, word
	case "empty":
		return TypeEmpty, word
	case "return":
		return Return, word
	case "if":
		return If, word
	case "true":
		return True, word
	case "false":
		return False, word
	case "bool":
		return TypeBool, word
	}

	return Identifier, word
}
