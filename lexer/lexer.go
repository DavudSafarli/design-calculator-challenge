package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Token interface{}
type MatcherFunc func(*Lexer) (Token Token, found bool)

type Lexer struct {
	tokens   []int
	matchers map[int]MatcherFunc

	r   *bufio.Reader
	err error
}

type Options struct {
	Tokens   []int
	Matchers map[int]MatcherFunc
}

func NewLexer(ops Options) Lexer {
	return Lexer{
		tokens:   ops.Tokens,
		matchers: ops.Matchers,
	}
}

type UnknownSymbolError struct {
	Symbol rune
}

func (e UnknownSymbolError) Error() string {
	return fmt.Sprintf("unknown character %q", e.Symbol)
}

func (l *Lexer) Lex(input string) ([]Token, error) {

	r := strings.NewReader(input)
	br := bufio.NewReader(r)
	l.r = br
	l.err = nil
	tokens := []Token{}

OUTER:
	for {
		if l.err == io.EOF {
			break
		}
		if l.err != nil {
			return nil, l.err
		}
		var token Token
		var found bool
		for _, tokenType := range l.tokens {
			matcherFn, ok := l.matchers[tokenType]
			if !ok {
				panic("no matcher exists for ..")
			}
			token, found = matcherFn(l)
			if found {
				tokens = append(tokens, token)
			}
			if l.err != nil {
				break OUTER
			}

			if found {
				break
			}
		}

		if !found {
			ch, _ := l.ReadNext()
			return nil, UnknownSymbolError{ch}
		}
	}
	return tokens, nil
}

func (l *Lexer) ReadNext() (rune, error) {
	ch, _, err := l.r.ReadRune()
	if err != nil && err != io.EOF {
		panic(err)
	}
	l.err = err
	return ch, err
}

func (l *Lexer) Unread() {
	if l.err == io.EOF {
		l.err = nil
	}
	err := l.r.UnreadRune()
	l.err = err
}

func (l *Lexer) ReadInt() (string, bool) {
	str, ok := l.ReadBetween('0', '9')
	return str, ok
}

func (l *Lexer) ReadIntOrFloat() (string, bool) {
	decimal, ok := l.ReadBetween('0', '9')
	if !ok {
		return decimal, ok
	}
	if !l.ReadChar('.') {
		return decimal, ok
	}
	floating, ok := l.ReadBetween('0', '9')
	// TODO: can't unread more than 1 times. One options is keep a customizable-sized Queue for read Runes.
	// For now, treat cases like "123." as "123.0" and take them as valid floats
	if !ok {
		return decimal + ".", true
	}
	return decimal + "." + floating, ok
}

func (l *Lexer) ReadChar(want rune) bool {
	ch, err := l.ReadNext()
	if err != nil {
		return false
	}
	if ch == want {
		return true
	}
	l.Unread()
	return false
}

func (l *Lexer) ReadBetween(from, to rune) (string, bool) {
	ch, err := l.ReadNext()
	if err != nil {
		return "", false
	}

	if !isBetween(ch, from, to) {
		l.Unread()
		return "", false
	}

	sb := strings.Builder{}
	for isBetween(ch, from, to) {
		sb.WriteRune(ch)
		ch, err = l.ReadNext()
		if err != nil {
			return sb.String(), true
		}
	}
	// no error, no io.EOF
	l.Unread()
	return sb.String(), true
}

func (l *Lexer) ReadUntil(runes []rune) (string, bool) {
	ch, err := l.ReadNext()
	if err != nil {
		return "", false
	}

	isOneOf := func(ch rune) bool {
		for _, v := range runes {
			if ch == v {
				return true
			}
		}
		return false
	}
	if !isOneOf(ch) {
		l.Unread()
		return "", false
	}

	sb := strings.Builder{}
	for isOneOf(ch) {
		sb.WriteRune(ch)
		ch, err = l.ReadNext()
		if err != nil {
			return sb.String(), true
		}
	}
	// no error, no io.EOF
	l.Unread()
	return sb.String(), true
}

func isBetween(ch, from, to rune) bool {
	return ch >= from && ch <= to
}
