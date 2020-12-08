package memorymodel

import "fmt"

type ContextStack struct {
	stack []*Context
}

func NewContextStack() *ContextStack {
	return &ContextStack{
		stack: []*Context{},
	}
}

func (cs *ContextStack) Push(context *Context) {
	cs.stack = append(cs.stack, context)
}

func (cs *ContextStack) Peek() *Context {
	if cs.Size() > 0 {
		return cs.stack[cs.Size()-1]
	}
	return nil
}

func (cs *ContextStack) Pop() error {
	if cs.Size() > 0 {
		cs.stack = cs.stack[:cs.Size()-1]
	}
	return fmt.Errorf("stack is empty")
}

func (cs *ContextStack) Size() int {
	return len(cs.stack)
}
