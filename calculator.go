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
	var operators TokenStack

	// addOperandNode takes 2 node from stack, creates a new Expression Node like below,
	//     +
	//   /   \
	//  a     b
	// and pushes it back to the slice of postfix nodes
	addOperandNode := func(op Token) {
		b, _ := postfix.Pop()
		a, _ := postfix.Pop()
		if op.IsAddOP() {
			postfix.Push(AddNode{a, b})
		}
		if op.IsSubOP() {
			postfix.Push(SubNode{a, b})
		}
		if op.IsMulOP() {
			postfix.Push(MulNode{a, b})
		}
		if op.IsDivOP() {
			postfix.Push(DivNode{a, b})
		}
		if op.IsPowOP() {
			postfix.Push(PowNode{a, b})
		}
	}
	prev := Token{}
	for _, token := range tokens {
		if token.IsNum() {
			val, _ := strconv.ParseFloat(token.Value, 64)
			postfix.Push(NumNode{val})
		} else if token.IsOP() {
			for {
				prevOP, exists := operators.Top()
				if !exists || prevOP.Type == L_PAR {
					break
				}
				if precedence[prevOP.Type] < precedence[token.Type] {
					break
				}
				operators.Pop() // remove element
				addOperandNode(prevOP)

			}
			operators.Push(token)
		} else if token.IsLeftParacentesis() {
			// if there is a NUM before LEFT_PAR, then consider it as multiplication
			if prev.IsNum() {
				operators.Push(Token{MUL, ""})
			}
			operators.Push(token)
		} else if token.IsRightParacentesis() {
			for {
				prevOP, _ := operators.Pop()
				if prevOP.Type == L_PAR {
					break
				}
				addOperandNode(prevOP)
			}
		}
		prev = token
	}

	for !operators.IsEmpty() {
		op, _ := operators.Pop()
		addOperandNode(op)
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
