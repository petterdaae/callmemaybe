package grammar

import (
	"fmt"
	"strconv"
)

const (
	rax = "rax"
	rbx = "rbx"
	rcx = "rxc"
	rsp = "rsp" // stack pointer
	add = "add"
)

type AssemblyOutput struct {
	Operations []string
	StackSize int
	Identifiers map[string]int
}

func (output *AssemblyOutput) Start() {
	output.Operations = append(output.Operations, "extern printf")
	output.Operations = append(output.Operations, "global main")
	output.Operations = append(output.Operations, "section .data")
	output.Operations = append(output.Operations, "format: db '%x', 10, 0")
	output.Operations = append(output.Operations, "section .text")
	output.Operations = append(output.Operations, "main:")
}

func (output *AssemblyOutput) End() {
	output.Operations = append(output.Operations, "mov rax, 60")
	output.Operations = append(output.Operations, "syscall")
}

func (output *AssemblyOutput) move(destination string, source string) {
	line := fmt.Sprintf("mov %s, %s", destination, source)
	output.Operations = append(output.Operations, line)
}

func (output *AssemblyOutput) add(destination string, value string) {
	line := fmt.Sprintf("%s %s, %s", add, destination, value)
	output.Operations = append(output.Operations, line)
}

func (output *AssemblyOutput) mult(destination string, value string) {
	line := fmt.Sprintf("imul %s, %s", destination, value)
	output.Operations = append(output.Operations, line)
}

func (output *AssemblyOutput) push(value string) {
	line := fmt.Sprintf("push %s", value)
	output.StackSize++
	output.Operations = append(output.Operations, line)
}

func (output *AssemblyOutput) pop(destination string) {
	line := fmt.Sprintf("pop %s", destination)
	output.StackSize--
	output.Operations = append(output.Operations, line)
}

func (output *AssemblyOutput) println(value string) {
	output.Operations = append(output.Operations, fmt.Sprintf("mov rdi, format"))
	output.Operations = append(output.Operations, fmt.Sprintf("mov rsi, %s", value))
	output.Operations = append(output.Operations, fmt.Sprintf("xor rax, rax"))
	output.Operations = append(output.Operations, fmt.Sprintf("call printf"))
}

func (exp ExpPlus) Generate(output *AssemblyOutput) error {
	err := exp.Left.Generate(output)
	if err != nil {
		return fmt.Errorf("failed to generate code for left side of plus exp: %w", err)
	}
	output.push(rax)
	err = exp.Right.Generate(output)
	if err != nil {
		return fmt.Errorf("failed to generate code for right side of plus exp: %w", err)
	}
	output.pop(rbx)
	output.add(rax, rbx)
	return nil
}

func (exp ExpMultiply) Generate(output *AssemblyOutput) error {
	err := exp.Left.Generate(output)
	if err != nil {
		return fmt.Errorf("failed to generate code for left side of multiply exp: %w", err)
	}
	output.push(rax)
	err = exp.Right.Generate(output)
	if err != nil {
		return fmt.Errorf("failed to generate code for right side of multiply exp: %w ", err)
	}
	output.pop(rbx)
	output.mult(rax, rbx)
	return nil
}

func (exp ExpParentheses) Generate(output *AssemblyOutput) error {
	err := exp.Inside.Generate(output)
	if err != nil {
		return fmt.Errorf("failed to generate code for inside of parentheses exp: %w", err)
	}
	return nil
}

func (exp ExpNum) Generate(output *AssemblyOutput) error {
	val := strconv.Itoa(exp.Value)
	output.move(rax, fmt.Sprintf("0x%s", val))
	return nil
}

func (exp ExpLet) Generate(output *AssemblyOutput) error {
	err := exp.IdentifierExp.Generate(output)
	if err != nil {
		return fmt.Errorf("failed to generate code for identifier exp in let exp: %w", err)
	}
	output.push(rax)
	output.Identifiers[exp.Identifier] = output.StackSize
	err = exp.Inside.Generate(output)
	if err != nil {
		return fmt.Errorf("failed to generate code for exp inside let exp: %w", err)
	}
	output.pop(rbx)
	return nil
}

func (exp ExpIdentifier) Generate(output *AssemblyOutput) error {
	identifierStackPos, ok := output.Identifiers[exp.Name]
	if !ok {
		return fmt.Errorf("uknown identifier: %s", exp.Name)
	}
	diff := (output.StackSize - identifierStackPos) * 8
	if diff < 0 {
		return fmt.Errorf("negative stack position for identifier: %s", exp.Name)
	}
	identifierAddr := fmt.Sprintf("[%s+%d]", rsp, diff)
	output.move(rax, identifierAddr)
	return nil
}

func (stmt StmtSeq) Generate(output *AssemblyOutput) error {
	for i := range stmt.Statements {
		err := stmt.Statements[i].Generate(output)
		if err != nil {
			return fmt.Errorf("failed to generate code for statement in sequence: %w", err)
		}
	}
	return nil
}

func (stmt StmtAssign) Generate(output *AssemblyOutput) error {
	err := stmt.Expression.Generate(output)
	if err != nil {
		return fmt.Errorf("failed to generate code for expression in assign statement: %w", err)
	}
	output.push(rax)
	output.Identifiers[stmt.Identifier] = output.StackSize
	return nil
}

func (stmt StmtPrintln) Generate(output *AssemblyOutput) error {
	err := stmt.Expression.Generate(output)
	if err != nil {
		return fmt.Errorf("failed to generate code for expression in println: %w", err)
	}
	output.println(rax)
	return nil
}
