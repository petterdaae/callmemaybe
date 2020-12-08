package assemblyoutput

import "fmt"

type ProcedureStack struct {
	stack []*procedure
}

func NewProcedureStack() *ProcedureStack {
	return &ProcedureStack{
		stack: []*procedure{},
	}
}

func (ps *ProcedureStack) Push(procedure *procedure) {
	ps.stack = append(ps.stack, procedure)
}

func (ps *ProcedureStack) Peek() *procedure {
	if ps.Size() > 0 {
		return ps.stack[ps.Size()-1]
	}
	return nil
}

func (ps *ProcedureStack) Pop() error {
	if ps.Size() > 0 {
		ps.stack = ps.stack[:ps.Size()-1]
	}
	return fmt.Errorf("stack is empty")
}

func (ps *ProcedureStack) Size() int {
	return len(ps.stack)
}
