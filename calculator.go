package calculator

import (
	"fmt"
	"strconv"
	"strings"
)

// BODMAS
var precedence = map[int]int{
	POW: 3,
	MUL: 2,
	DIV: 2,
	ADD: 1,
	SUB: 1,
}

type Calculator struct {
	tokenizer Tokenizer
}

func NewCalculator() Calculator {
	return Calculator{
		tokenizer: Tokenizer{},
	}
}

func (c Calculator) Eval(input string) float64 {
	reader := strings.NewReader(input)

	tokens, _ := c.tokenizer.Tokenize(reader)
	headNode := c.buildExpressionTree(tokens)
	fmt.Println(headNode)
	return headNode.Calculate()
}

func (c Calculator) buildExpressionTree(tokens []Token) Calculatable {
	var postfix CalculatableStack
	var operands TokenStack

	// addOperandNode takes 2 node from stack, creates a new Expression Node like below,
	//     +
	//   /   \
	//  a     b
	// and pushes it back to the slice of postfix nodes
	addOperandNode := func(operand Token) {
		b, _ := postfix.Pop()
		a, _ := postfix.Pop()
		if operand.IsAddOperand() {
			postfix.Push(AddNode{a, b})
		}
		if operand.IsSubOperand() {
			postfix.Push(SubNode{a, b})
		}
		if operand.IsMulOperand() {
			postfix.Push(MulNode{a, b})
		}
		if operand.IsDivOperand() {
			postfix.Push(DivNode{a, b})
		}
		if operand.IsPowOperand() {
			postfix.Push(PowNode{a, b})
		}
	}

	for _, token := range tokens {
		if token.IsNum() {
			val, _ := strconv.ParseFloat(token.Value, 64)
			postfix.Push(NumNode{val})
			continue
		}
		if token.IsOperand() {
			for {
				prevOperand, ok := operands.Top()
				if ok && prevOperand.Type != L_PAR {
					if precedence[prevOperand.Type] >= precedence[token.Type] {
						operands.Pop() // remove element
						addOperandNode(prevOperand)
					} else {
						break
					}
				} else {
					break
				}
			}
			operands.Push(token)
			continue
		}
		if token.IsLeftParacentesis() {
			operands.Push(token)
			continue
		}
		if token.IsRightParacentesis() {
			for {
				prevOperand, _ := operands.Pop()
				if prevOperand.Type == L_PAR {
					break
				}
				addOperandNode(prevOperand)
			}
		}
	}

	for !operands.IsEmpty() {
		operand, _ := operands.Pop()
		addOperandNode(operand)
	}
	node, _ := postfix.Pop()
	return node
}
