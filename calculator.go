package calculator

import (
	"strconv"

	"github.com/DavudSafarli/design-calculator-challenge/lexer"
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
	lexer lexer.Lexer
}

func NewCalculator() Calculator {
	lexer := buildLexerWithBODMASSupport()
	return Calculator{
		lexer: lexer,
	}
}

func (c Calculator) Eval(input string) float64 {
	lexerTokens, _ := c.lexer.Lex(input)
	tokens := make([]Token, len(lexerTokens))

	for _, v := range lexerTokens {
		tokens = append(tokens, v.(Token))
	}
	headNode := c.buildExpressionTree(tokens)
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
				prevOperand, exists := operands.Top()
				if !exists || prevOperand.Type == L_PAR {
					break
				}
				if precedence[prevOperand.Type] < precedence[token.Type] {
					break
				}
				operands.Pop() // remove element
				addOperandNode(prevOperand)

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

func buildLexerWithBODMASSupport() lexer.Lexer {
	// 1-char matcher function for Lexer
	createOneCharMatcher := func(ch rune, tokenType int) lexer.MatcherFunc {
		return func(l *lexer.Lexer) (token lexer.Token, found bool) {
			r, _ := l.ReadNext()
			if r == ch {
				return Token{tokenType, ""}, true
			}
			l.Unread()
			return nil, false
		}
	}

	return lexer.NewLexer(lexer.Options{
		Tokens:      []int{NUM, ADD, SUB, MUL, DIV, POW, L_PAR, R_PAR},
		CharsToPass: []rune{' ', '\t', '\n'},
		Matchers: map[int]lexer.MatcherFunc{
			ADD:   createOneCharMatcher('+', ADD),
			SUB:   createOneCharMatcher('-', SUB),
			MUL:   createOneCharMatcher('*', MUL),
			DIV:   createOneCharMatcher('/', DIV),
			POW:   createOneCharMatcher('^', POW),
			L_PAR: createOneCharMatcher('(', L_PAR),
			R_PAR: createOneCharMatcher(')', R_PAR),
			NUM: func(l *lexer.Lexer) (token lexer.Token, found bool) {
				val, ok := l.ReadIntOrFloat()
				if !ok {
					return Token{}, false
				}
				return Token{NUM, val}, true
			},
		},
	})
}
