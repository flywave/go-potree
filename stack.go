package potree

import (
	"errors"
	"fmt"
)

type Stack struct {
	Element []*Node
}

func NewStack() *Stack {
	return &Stack{}
}

func (stack *Stack) Push(value ...*Node) {
	stack.Element = append(stack.Element, value...)
}

func (stack *Stack) Top() (value *Node) {
	if stack.Size() > 0 {
		return stack.Element[stack.Size()-1]
	}
	return nil
}

func (stack *Stack) Pop() (err error) {
	if stack.Size() > 0 {
		stack.Element = stack.Element[:stack.Size()-1]
		return nil
	}
	return errors.New("Stack is empty.")
}

func (stack *Stack) Swap(other *Stack) {
	switch {
	case stack.Size() == 0 && other.Size() == 0:
		return
	case other.Size() == 0:
		other.Element = stack.Element[:stack.Size()]
		stack.Element = nil
	case stack.Size() == 0:
		stack.Element = other.Element
		other.Element = nil
	default:
		stack.Element, other.Element = other.Element, stack.Element
	}
	return
}

func (stack *Stack) Set(idx int, value *Node) (err error) {
	if idx >= 0 && stack.Size() > 0 && stack.Size() > idx {
		stack.Element[idx] = value
		return nil
	}
	return errors.New("Set faile!")
}

func (stack *Stack) Get(idx int) (value *Node) {
	if idx >= 0 && stack.Size() > 0 && stack.Size() > idx {
		return stack.Element[idx]
	}
	return nil
}

func (stack *Stack) Size() int {
	return len(stack.Element)
}

func (stack *Stack) Empty() bool {
	if stack.Element == nil || stack.Size() == 0 {
		return true
	}
	return false
}

func (stack *Stack) Print() {
	for i := len(stack.Element) - 1; i >= 0; i-- {
		fmt.Println(i, "=>", stack.Element[i])
	}
}
