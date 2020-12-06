package assembly

import "fmt"

type ContextElementType int

const (
	Stack ContextElementType = iota
	Heap
	Procedure
	Invalid
)

type Contexts struct {
	stack []*context
}

type context struct {
	procedures               map[string]string
	stack                    map[string]int
	propagateChangesToParent bool
	procedure                *procedure
}

type procedure struct{
	name string
	operations []string
	stackSizeWhenCreated int
}

func (contexts Contexts) Push(context *context) {
	contexts.stack = append(contexts.stack, context)
}

func (contexts Contexts) Size() int {
	return len(contexts.stack)
}

func (contexts Contexts) Peek() *context {
	size := contexts.Size()
	if size == 0 {
		return nil
	}
	return contexts.stack[size-1]
}

func (contexts Contexts) Pop() *context {
	result := contexts.Peek()
	if result == nil {
		return nil
	}
	size := contexts.Size()
	contexts.stack = contexts.stack[:size-1]
	return result
}

func (contexts Contexts) getFromContext(context *context, name string, stackSize int) (ContextElementType, string, error) {
	if context == nil {
		return Invalid, "", fmt.Errorf("context stack is empty")
	}

	procedure := contexts.getParentProcedure()

	stack, ok := context.stack[name]
	if ok {
		diff := (stackSize - stack) * 8

		if procedure != nil {
			if stack > procedure.stackSizeWhenCreated {
				return Stack, fmt.Sprintf("[rsp+%d]", diff), nil
			}

			return Stack, fmt.Sprintf("[rsp+rcx+%d+8]", diff), nil
		}

		return Stack, fmt.Sprintf("[rsp+%d]", diff), nil
	}

	proc, ok := context.procedures[name]
	if ok {
		return Procedure, proc, nil
	}

	return Invalid, "", fmt.Errorf("could not get '%s' from context", name)
}

func (contexts Contexts) getParentProcedure() *procedure {
	for i := contexts.Size() - 1; i >= 0; i-- {
		current := contexts.stack[i].procedure
		if current != nil {
			return current
		}
	}
	return nil
}

func (contexts Contexts) Get(name string, stackSize int) (ContextElementType, string, error) {
	context := contexts.Peek()
	return contexts.getFromContext(context, name, stackSize)
}

func (contexts Contexts) StackInsert(name string, value string, stackSize int) ([]string, int, error) {
	var operations []string
	pushes := 0
	top := contexts.Peek()
	if top == nil {
		return nil, 0, fmt.Errorf("context stack is empty")
	}

	// If name is not in current stack
	_, ok := top.stack[name]
	if !ok {
		operations = append(operations, fmt.Sprintf("push %s", value))
		pushes++
		top.stack[name] = stackSize
		return operations, pushes, nil
	}

	// If name is in current stack
	if top.procedure == nil {
		for i := contexts.Size() - 1; i >= 0; i-- {
			current := contexts.stack[i]
			_, ok := current.stack[name]
			if !ok {
				break
			}

			_, address, _ := contexts.getFromContext(current, name, stackSize)
			operations = append(operations, fmt.Sprintf("mov %s, %s", address, value))

			if current.procedure != nil {
				break
			}
		}
	}

	return operations, pushes, nil
}

func (contexts Contexts) ProcInsert(name string, alias string) error {
	return nil
}
