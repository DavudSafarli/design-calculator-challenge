package calculator

import (
	"strconv"
	"strings"
)

type Calculator struct {
	tokenizer Tokenizer
}

func (c Calculator) Eval(input string) float64 {
	reader := strings.NewReader(input)

	tokens, _ := c.tokenizer.Tokenize(reader)
	headNode := c.buildExpressionTree(tokens)
	return headNode.Calculate()
}

func (c Calculator) buildExpressionTree(tokens []Token) Calculatable {
	var calculatables Stack
	var operands Stack

	for _, token := range tokens {
		if token.IsNum() {
			val, _ := strconv.ParseFloat(token.Value, 64)
			calculatables.Push(NumNode{val})
			continue
		}
		if token.IsAddOperand() {
			operands.Push(token)
			continue
		}
	}

	for !operands.IsEmpty() {
		token, _ := operands.Pop()
		operand := token.(Token)

		c1, _ := calculatables.Pop()
		c2, _ := calculatables.Pop()
		a := c1.(Calculatable)
		b := c2.(Calculatable)

		if operand.IsAddOperand() {
			calculatables.Push(AddNode{a, b})
		}
	}
	node, _ := calculatables.Pop()
	return node.(Calculatable)
}

func NewCalculator() Calculator {
	return Calculator{
		tokenizer: Tokenizer{},
	}
}
