package tokenizer

import (
	"fmt"
)

type ProgramInput = []rune
type Tokenizer = func(ProgramInput, int) (string, int, error)

func Character(char rune) Tokenizer {
	return func(program ProgramInput, index int) (string, int, error) {
		if index >= len(program) {
			return "", index, fmt.Errorf("index out bounds in character tokenizer")
		}
		if char == program[index] {
			return string(char), index + 1, nil
		}
		return "", 0, fmt.Errorf("character tokenizer failed to eat character")
	}
}

func Word(word string) Tokenizer {
	runes := []rune(word)
	return func(program ProgramInput, index int) (string, int, error) {
		if index + len(runes) > len(program) {
			return "", index, fmt.Errorf("word to long for the rest of the program in word tokenizer")
		}
		for _, rune := range runes {
			eatRune := Character(rune)
			_, _, err := eatRune(program, index)
			if err != nil {
				return "", index, fmt.Errorf("character tokenizer failed in word tokenizer: %w", err)
			}
		}
		result := string(runes)
		return result, index + len(runes), nil
	}
}

func Num(program string, index int) (string, int, error) {
	return "", 0, nil
}
