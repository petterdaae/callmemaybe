package language

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
)

const (
	Number Token = iota
	Plus
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
	TypeEmpty
	Whitespace
	PrintLn
	Return
	Identifier
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
	}
	
	return Error, ""
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
	case "empty":
		return TypeEmpty, word
	case "return":
		return Return, word
	}

	return Identifier, word
}
