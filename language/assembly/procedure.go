package assembly

import (
	"fmt"
)

type Procedure struct {
	Name                     string
	Operations               []string
	StackSizeWhenInitialized int
	NumberOfArgs             int
	ReturnType               string
}

func (proc *Procedure) start(stackSize int) {
	proc.Operations = append(proc.Operations, fmt.Sprintf("mov rdx, %d", stackSize))
	proc.Operations = append(proc.Operations, fmt.Sprintf("sub rcx, rdx"))
	proc.Operations = append(proc.Operations, fmt.Sprintf("imul rcx, 8"))
}

func (proc *Procedure) end(stackSize int) int {
	diff := stackSize - proc.StackSizeWhenInitialized
	for i := 0; i < diff; i++ {
		proc.Operations = append(proc.Operations, fmt.Sprintf("pop rax"))
	}
	proc.Operations = append(proc.Operations, fmt.Sprintf("ret"))
	return diff
}
