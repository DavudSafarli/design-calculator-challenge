package lexer_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/DavudSafarli/design-calculator-challenge/lexer"
)

const (
	NUM = iota
	ADD
	SUB
	MUL
	DIV
	POW
	L_PAR
	R_PAR

	SPACE
)

func createOneCharMatcher(ch rune, tokenType int) lexer.MatcherFunc {
	return func(l *lexer.Lexer) (token lexer.Token, found bool) {
		ok := l.ReadChar(ch)
		if !ok {
			return Token{}, false
		}
		return Token{tokenType, ""}, true
	}
}

type Token struct {
	Type  int
	Value string
}

func TestLexer(t *testing.T) {
	lex := lexer.NewLexer(lexer.Options{
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
				val, ok := l.ReadUntil([]rune{' ', '\t', '\n'})
				if !ok {
					return Token{}, false
				}
				return Token{SPACE, val}, true
			},
		},
	})

	tests := []struct {
		input string
		want  []lexer.Token
	}{
		{"", []lexer.Token{}},
		{"123456789", []lexer.Token{Token{NUM, "123456789"}}},
		{"  1  ", []lexer.Token{Token{SPACE, "  "}, Token{NUM, "1"}, Token{SPACE, "  "}}},
		{"+-*/^()", []lexer.Token{Token{ADD, ""}, Token{SUB, ""}, Token{MUL, ""}, Token{DIV, ""}, Token{POW, ""}, Token{L_PAR, ""}, Token{R_PAR, ""}}},
		{
			"7-7/7+7*14",
			[]lexer.Token{Token{NUM, "7"}, Token{SUB, ""}, Token{NUM, "7"}, Token{DIV, ""}, Token{NUM, "7"}, Token{ADD, ""}, Token{NUM, "7"}, Token{MUL, ""}, Token{NUM, "14"}},
		},
		{"123.456+5.", []lexer.Token{Token{NUM, "123.456"}, Token{ADD, ""}, Token{NUM, "5."}}},
		{"2.25*4", []lexer.Token{Token{NUM, "2.25"}, Token{MUL, ""}, Token{NUM, "4"}}},
		{
			"(2+2.22)",
			[]lexer.Token{Token{L_PAR, ""}, Token{NUM, "2"}, Token{ADD, ""}, Token{NUM, "2.22"}, Token{R_PAR, ""}},
		},
	}

	for _, tt := range tests {
		testName := fmt.Sprint("Tokenizing", tt.input)
		t.Run(testName, func(t *testing.T) {
			got, err := lex.Lex(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokenize() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestLexeInvalidInput(t *testing.T) {
	lex := lexer.NewLexer(lexer.Options{
		Tokens: []int{ADD},
		Matchers: map[int]lexer.MatcherFunc{
			ADD: createOneCharMatcher('+', ADD),
		},
	})

	_, err := lex.Lex("+-")
	expected := lexer.UnknownSymbolError{'-'}

	if err != expected {
		t.Fatalf("\nexpected: %v\nactual  : %v", expected, err)
	}
}

func TestLexeIncorrectUsage(t *testing.T) {
	lex := lexer.NewLexer(lexer.Options{
		Tokens: []int{ADD},
		Matchers: map[int]lexer.MatcherFunc{
			ADD: func(l *lexer.Lexer) (lexer.Token, bool) {
				r, done := l.ReadNext()
				if done {
					return Token{}, false
				}
				if r == '+' {
					return Token{ADD, ""}, true
				}
				// l.Unread() //
				return Token{}, false
			},
		},
	})

	_, err := lex.Lex("+-")

	expected := lexer.ErrMatcherForgotToUnread
	if err != expected {
		t.Fatalf("\nexpected: %v\nactual  : %v", expected, err)
	}

	if err != expected {
		t.Fatal(err)
	}
}
