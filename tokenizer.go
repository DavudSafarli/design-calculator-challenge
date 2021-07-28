package calculator

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Operators include:
const (
	OP_ADD = '+'
	OP_SUB = '-'
	OP_MUL = '*'
	OP_DIV = '/'
)

// Token Types include:
// NUM - can be an integer(1) or float(1.2)
// ADD - addition operand (+)
// SUB - addition operand (-)
// MUL - addition operand (*)
// DIV - addition operand (/)
const (
	NUM = iota
	ADD
	SUB
	MUL
	DIV
)

type Token struct {
	Type  int
	Value string
}

func (t Token) IsNum() bool {
	return t.Type == NUM
}

func (t Token) IsOperand() bool {
	return t.Type == ADD || t.Type == SUB || t.Type == MUL || t.Type == DIV
}

func (t Token) IsAddOperand() bool {
	return t.Type == ADD
}
func (t Token) IsSubOperand() bool {
	return t.Type == SUB
}
func (t Token) IsMulOperand() bool {
	return t.Type == MUL
}
func (t Token) IsDivOperand() bool {
	return t.Type == DIV
}

type Tokenizer struct {
	r *bufio.Reader
}

func (t *Tokenizer) readRune() (rune, error) {
	r, _, err := t.r.ReadRune()
	return r, err
}
func (t *Tokenizer) unread() {
	err := t.r.UnreadRune()
	if err != nil {
		panic(fmt.Errorf("error while unreading: %w", err))
	}
}

// Tokenize extracts all the tokens from the input, and returns a slice of Tokens
func (t *Tokenizer) Tokenize(r io.Reader) ([]Token, error) {
	t.r = bufio.NewReader(r)
	tokens := []Token{}
	for {
		token, err := t.getNextToken()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

// getNextToken finds the next token(excluding the whitespaces)
func (t *Tokenizer) getNextToken() (Token, error) {
	ch, err := t.readRune()
	if err != nil {
		return Token{}, err
	}
	t.unread()

	if isWhiteSpace(ch) {
		err := t.bypassWhitespaces()
		if err != nil {
			return Token{}, err
		}
		ch, err = t.readRune()
		if err != nil {
			return Token{}, err
		}
		t.unread()
	}
	if isNumeric(ch) {
		token, err := t.ReadNum()
		return token, err
	}
	if isOperator(ch) {
		token, err := t.ReadOperator()
		return token, err
	}

	// TODO: give more details about where the error was
	return Token{}, fmt.Errorf("invalid input")
}

func isWhiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}
func isNumeric(ch rune) bool {
	return ch >= '0' && ch <= '9'
}
func isOperator(ch rune) bool {
	return ch == OP_ADD || ch == OP_SUB || ch == OP_MUL || ch == OP_DIV
}

// bypassWhitespaces reads and passes all contiguous white spaces
func (t *Tokenizer) bypassWhitespaces() error {
	ch, err := t.readRune()
	if err != nil {
		return err
	}
	for isWhiteSpace(ch) {
		ch, err = t.readRune()
		if err != nil {
			return err
		}
	}
	t.unread()
	return nil
}

// ReadNum reads a NUM token(integer or float). returns error if fails to find NUM-type Token
func (t *Tokenizer) ReadNum() (token Token, err error) {
	decimal, err := t.ReadInt()
	if err == io.EOF {
		return Token{NUM, decimal}, nil
	}
	if err != nil {
		return Token{}, nil
	}

	ch, err := t.readRune()
	if err != nil {
		return Token{}, err
	}
	if ch == '.' {
		floating, err := t.ReadInt()
		if floating == "" {
			return Token{}, fmt.Errorf("expecting at least 1 number after dot")
		}
		value := fmt.Sprintf("%v.%v", decimal, floating)
		if err == io.EOF {
			return Token{NUM, value}, nil
		}
		if err != nil {
			return Token{}, err
		}
		t.unread()
		return Token{NUM, value}, nil
	}

	t.unread()
	token = Token{NUM, decimal}
	return token, nil
}

// read Integer
func (t *Tokenizer) ReadInt() (string, error) {
	ch, err := t.readRune()
	if err != nil {
		return "", err
	}

	if !isNumeric(ch) {
		return "", fmt.Errorf("expecting NUM Token, found non-numeric")
	}

	sb := strings.Builder{}
	for isNumeric(ch) {
		sb.WriteRune(ch)
		ch, err = t.readRune()
		if err != nil {
			return sb.String(), err
		}
	}
	// no error, no io.EOF
	t.unread()
	return sb.String(), err
}

// ReadOperator reads an Operator-type token(ADD, etc.). return error if fails to find any Operator type Token
func (t *Tokenizer) ReadOperator() (Token, error) {
	ch, err := t.readRune()
	if err != nil {
		return Token{}, err
	}

	if ch == OP_ADD {
		return Token{ADD, ""}, nil
	}
	if ch == OP_SUB {
		return Token{SUB, ""}, nil
	}
	if ch == OP_MUL {
		return Token{MUL, ""}, nil
	}
	if ch == OP_DIV {
		return Token{DIV, ""}, nil
	}
	return Token{}, fmt.Errorf("expecting Operator-type Token, but can't find")
}
