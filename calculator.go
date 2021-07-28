package calculator

import (
	"fmt"
	"strconv"
	"strings"
)

// BODMAS
var precedence = map[int]int{
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
	var calculatables CalculatableStack
	var operands TokenStack

	addOperandNode := func(operand Token) {
		b, _ := calculatables.Pop()
		a, _ := calculatables.Pop()
		if operand.IsAddOperand() {
			calculatables.Push(AddNode{a, b})
		}
		if operand.IsSubOperand() {
			calculatables.Push(SubNode{a, b})
		}
		if operand.IsMulOperand() {
			calculatables.Push(MulNode{a, b})
		}
		if operand.IsDivOperand() {
			calculatables.Push(DivNode{a, b})
		}
	}

	for _, token := range tokens {
		if token.IsNum() {
			val, _ := strconv.ParseFloat(token.Value, 64)
			calculatables.Push(NumNode{val})
			continue
		}
		if token.IsOperand() {
			for {
				prevOperand, ok := operands.Top()
				if ok {
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
	}

	for !operands.IsEmpty() {
		operand, _ := operands.Pop()
		addOperandNode(operand)
	}
	node, _ := calculatables.Pop()
	return node
}
