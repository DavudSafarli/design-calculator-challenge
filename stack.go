package calculator

type TokenStack []Token

// IsEmpty: check if stack is empty
func (s *TokenStack) IsEmpty() bool {
	return len(*s) == 0
}

// Push a new value onto the stack
func (s *TokenStack) Push(str Token) {
	*s = append(*s, str) // Simply append the new value to the end of the stack
}

// Remove and return top element of stack. Return false if stack is empty.
func (s *TokenStack) Pop() (Token, bool) {
	if s.IsEmpty() {
		return Token{}, false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]      // Remove it from the stack by slicing it off.
		return element, true
	}
}

// Return top element of stack. Return false if stack is empty.
func (s *TokenStack) Top() (Token, bool) {
	if s.IsEmpty() {
		return Token{}, false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		return element, true
	}
}

type CalculatableStack []Calculatable

// IsEmpty: check if stack is empty
func (s *CalculatableStack) IsEmpty() bool {
	return len(*s) == 0
}

// Push a new value onto the stack
func (s *CalculatableStack) Push(str Calculatable) {
	*s = append(*s, str) // Simply append the new value to the end of the stack
}

// Remove and return top element of stack. Return false if stack is empty.
func (s *CalculatableStack) Pop() (Calculatable, bool) {
	if s.IsEmpty() {
		return nil, false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]      // Remove it from the stack by slicing it off.
		return element, true
	}
}

// Return top element of stack. Return false if stack is empty.
func (s *CalculatableStack) Top() (Calculatable, bool) {
	if s.IsEmpty() {
		return nil, false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		return element, true
	}
}
