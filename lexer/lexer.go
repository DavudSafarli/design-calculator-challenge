package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Token represents any meaningful data extracted during the process of Lexical analysis.
type Token interface{}

// MatcherFunc is any user-defined function that tries to parse the input, extrant and return a Token.
// Any type of Token(defined by user) can be created and returned.
type MatcherFunc func(*Lexer) (Token Token, found bool)

// Lexer takes the grammar rules and uses them to tokenize* the input
type Lexer struct {
	tokens   []int
	matchers map[int]MatcherFunc

	r   *bufio.Reader
	err error
}

// Options represents the grammar rules.
//
// Tokens is a slice of unique TokenTypes to identify the Tokens.
//
// Matchers holds rules to find Tokens in the text. Keys of the map are the TokenTypes passed to `Tokens` options
type Options struct {
	Tokens   []int
	Matchers map[int]MatcherFunc
}

// NewLexer creates a new Lexer with given options
func NewLexer(ops Options) Lexer {
	return Lexer{
		tokens:   ops.Tokens,
		matchers: ops.Matchers,
	}
}

// UnknownSymbolError is an error that Lexer returns when encountering a symbol that she can not recignize
type UnknownSymbolError struct {
	Symbol rune
}

func (e UnknownSymbolError) Error() string {
	return fmt.Sprintf("unknown character %q", e.Symbol)
}

// Lex starts the lexical analysis and returns the slice of Tokens
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

// ReadNext returns the next rune of the input
//
// After using ReadNext and not finding what you want,
// user is responsible for calling #Unread to not let other matchers pass that character
func (l *Lexer) ReadNext() (rune, error) {
	ch, _, err := l.r.ReadRune()
	if err != nil && err != io.EOF {
		panic(err)
	}
	l.err = err
	return ch, err
}

// Unread unreads the last rune. Cannot be called more than one
func (l *Lexer) Unread() {
	if l.err == io.EOF {
		l.err = nil
	}
	err := l.r.UnreadRune()
	l.err = err
}

// ReadInt tries to read an integer (\d+) if. returns the number in string format if found
func (l *Lexer) ReadInt() (string, bool) {
	str, ok := l.ReadBetween('0', '9')
	return str, ok
}

// ReadIntOrFloat tries to read an integer or float (\d+\.\d+). returns the number in string format if found
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

// ReadChar tries to read the requeted char.
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

// ReadBetween reads all next contiguous chars that are between [from, to]
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

// ReadUntil reads all next contiguous chars that are in the given `runes` slice
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
