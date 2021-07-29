package calculator

import (
	"errors"
	"fmt"
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

// Calculator evaluates the given arithmetic expression.
type Calculator struct {
	lexer lexer.Lexer
}

func New() Calculator {
	lexer := buildLexerWithBODMASSupport()
	return Calculator{
		lexer: lexer,
	}
}

// EvalError gives info about what was the error, where it starts and where it ends in the user input
// example:
// user input: " 1+*"
// err: EvalError{Err2Operators, 3, 4}
// (3, 4) represents the position of "*" which can be used to visualize the error in the front-facing application
type EvalError struct {
	Err      error
	StartPos int
	EndPos   int
}

func (e EvalError) Error() string {
	return fmt.Sprintf("%s in position (%v, %v)", e.Err.Error(), e.StartPos, e.EndPos)
}

// Eval calculates given mathematical expression and returns the result
func (c Calculator) Eval(input string) (float64, error) {
	return c.eval(input)
}

// eval calculates the expression by going through following steps:
// - performs lexical analysis which generates sequence of Tokens.
// - validate the expression, report the error and invalid index position.
// - parses the Tokens and builds an expression tree, where each node is a `Calculatable`.
// - running the calculation process starting from the head node of the tree and getting the result
func (c Calculator) eval(input string) (float64, error) {
	lexerTokens, err := c.lexer.Lex(input)
	if err != nil {
		return 0, err
	}
	tokens := make([]Token, 0, len(lexerTokens))

	// convert `lexer.Token`s to original `Token`s defined by us
	for _, v := range lexerTokens {
		tokens = append(tokens, v.(Token))
	}

	invalidTokenPos, err := validateExpression(tokens)
	if err != nil {
		if invalidTokenPos != -1 {
			startPos, endPos := findTokenPositionInRawInput(tokens, invalidTokenPos)
			return 0, EvalError{err, startPos, endPos}
		}
		return 0, EvalError{err, -1, -1}
	}

	headNode := c.buildExpressionTree(tokens)
	return headNode.Calculate(), nil
}

// buildExpressionTree creates the expression tree and returns the head node
// given the valid Infix slice of Tokens
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
		if token.IsSpace() {
			continue
		}
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

// buildLexerWithBODMASSupport creates and returns a Lexer with lexical support for BODMAS, and SPACE
func buildLexerWithBODMASSupport() lexer.Lexer {
	// 1-char matcher function for Lexer
	createOneCharMatcher := func(ch rune, tokenType int) lexer.MatcherFunc {
		return func(l *lexer.Lexer) (token lexer.Token, found bool) {
			r, _ := l.ReadNext()
			if r == ch {
				return Token{tokenType, string(ch)}, true
			}
			l.Unread()
			return nil, false
		}
	}
	spaces := []rune{' ', '\t', '\n'}
	return lexer.NewLexer(lexer.Options{
		Tokens: []int{NUM, ADD, SUB, MUL, DIV, POW, L_PAR, R_PAR, SPACE},
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
			SPACE: func(l *lexer.Lexer) (token lexer.Token, found bool) {
				val, ok := l.ReadUntil(spaces)
				if !ok {
					return Token{}, false
				}
				return Token{SPACE, val}, true
			},
		},
	})
}

var ErrOperationBeforeRightParacentesis = errors.New("cannot have an operation before a closing-paracentesis")
var ErrInconsistentParacentesisCount = errors.New("inconsistent paracentesis count")
var Err2Operators = errors.New("cannot have 2 operators side by side")
var ErrOperationAfterLeftParacantesis = errors.New("cannot have an operation after an opening-paracentesis")
var ErrCannotStartWithOperator = errors.New("expression cannot start with an operator")

// validateExpression checks if expression is valid. returns the invalid index of the Token
// -1 means that, even though there was an error, position cannot be found
func validateExpression(tokens []Token) (int, error) {
	openParCount := 0
	prev := Token{}
	for i, token := range tokens {
		if token.IsLeftParacentesis() {
			openParCount++
		}
		// case: "..1+)"
		if token.IsRightParacentesis() && prev.IsOP() {
			return i, ErrOperationBeforeRightParacentesis
		}
		// case: "(3))"
		if token.IsRightParacentesis() && openParCount == 0 {
			return i, ErrInconsistentParacentesisCount
		}
		if token.IsRightParacentesis() {
			openParCount--
		}

		// case: "3/*4"
		if token.IsOP() && prev.IsOP() {
			return i, Err2Operators
		}
		// case: "3(+"
		if token.IsOP() && prev.IsLeftParacentesis() {
			return i, ErrOperationAfterLeftParacantesis
		}
		// case: "*5"
		if token.IsOP() && i == 0 {
			return i, ErrCannotStartWithOperator
		}

		prev = token
	}

	// (5+4
	if openParCount != 0 {
		return -1, ErrInconsistentParacentesisCount
	}
	return 0, nil
}

// findTokenPositionInRawInput find and returns the position of the Token in the Raw input
func findTokenPositionInRawInput(tokens []Token, tokenPos int) (startPos, endPos int) {
	for i := 0; i < tokenPos; i++ {
		startPos += len(tokens[i].Value)
	}
	endPos = startPos + len(tokens[tokenPos].Value)
	return startPos, endPos
}
