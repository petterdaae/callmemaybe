package tokenizer

import (
	"strings"
	"testing"
)

func TestIdentifier(t *testing.T) {
	tokenizer := New(strings.NewReader("foo=bar"))
	first, _ := tokenizer.NextToken()
	second, _ := tokenizer.NextToken()
	third, _ := tokenizer.NextToken()

	if first != Identifier || second != Assign || third != Identifier {
		t.Error()
	}
}

func TestLet(t *testing.T) {
	tokenizer := New(strings.NewReader("let+in"))
	first, _ := tokenizer.NextToken()
	tokenizer.NextToken()
	third, _ := tokenizer.NextToken()

	if first != Let || third != In {
		t.Error()
	}
}