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

func (contexts Contexts) Get(name string, stackSize int, procedure Procedure) (string, error) {
	top := contexts.Peek()
	if top == nil {
		return "", fmt.Errorf("context stack is empty")
	}

	stack, ok := top.stack[name]
	if ok {

	}
}
