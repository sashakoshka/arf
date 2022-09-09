package types

// Stack is a generic stack data structure.
type Stack[VALUE any] []VALUE

// Push adds a value to the top of the stack.
func (stack *Stack[VALUE]) Push (value VALUE) {
	*stack = append(*stack, value)
}

// Pop removes the topmost value of the stack and returns it.
func (stack *Stack[VALUE]) Pop () (value VALUE) {
	if len(*stack) < 1 {
		panic("can't pop off of empty stack")
	}

	newLength := len(*stack) - 1
	value  = (*stack)[newLength]
	*stack = (*stack)[:newLength]
	return
}

// Top returns the value on top of the stack.
func (stack Stack[VALUE]) Top () (value VALUE) {
	if len(stack) < 1 {
		panic("can't get top of empty stack")
	}
	
	value = stack[len(stack) - 1]
	return
}
