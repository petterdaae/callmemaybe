package tokenizer

import "testing"

func TestCharacter(t *testing.T) {
	ct1 := Character('a')
	ct2 := Character('b')
	ct3 := Character('c')

	program := []rune("abcd")
	index := 0
	var err error
	_, index, err = ct1(program, index)
	if err != nil {
		t.Error()
	}
	_, index, err = ct2(program, index)
	if err != nil {
		t.Error()
	}
	_, index, err = ct3(program, index)
	if err != nil {
		t.Error()
	}
	_, index, err = ct3(program, index)
	if err == nil {
		t.Error()
	}
}
