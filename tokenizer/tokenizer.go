package tokenizer

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
	ParenthesesStart
	ParenthesesEnd
	Whitespace
	EOF
	Error
)

const eof = rune(0)

type Token int

type Tokenizer struct {
	reader *bufio.Reader
}

func New(reader io.Reader) *Tokenizer {
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

	switch character {
	case eof:
		return EOF, ""
	case '*':
		return Multiply, string(character)
	case '+':
		return Plus, string(character)
	case '(':
		return ParenthesesStart, string(character)
	case ')':
		return ParenthesesEnd, string(character)
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
