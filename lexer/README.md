# Rune base Lexer

This package allows you to define `TokenMatchers` to extract token out of the input text.
MatcherFunctions can use `ReadNext` to get the next rune, check if matches the wanted Token pattern. If yes, you can simply return the user-defined Token type. If not, user should pay attention to `Unread` to notify Lexer to go 1 step back.

Simple example of Lexer with ADD(+) and SUB(-):
```go
func main() {
	type Token struct {
		Type  int
		Value string
	}

	TOKEN_ADD := 1
	TOKEN_SUB := -1

	oneCharMatcher := func(ch rune, tokenType int, name string) lexer.MatcherFunc {
		return func(l *lexer.Lexer) (token lexer.Token, found bool) {
			ok := l.ReadChar(ch)
			if !ok {
				return Token{}, false
			}
			return Token{tokenType, name}, true
		}
	}

	lex := lexer.NewLexer(lexer.Options{
		Tokens: []int{TOKEN_ADD, TOKEN_SUB},
		Matchers: map[int]lexer.MatcherFunc{
			TOKEN_ADD: oneCharMatcher('+', TOKEN_ADD, "ADD"),
			TOKEN_SUB: oneCharMatcher('-', TOKEN_SUB, "SUB"),
		},
	})

	lexerTokens, _ := lex.Lex("++--")

	fmt.Println(lexerTokens)
	// [{1 ADD} {1 ADD} {-1 SUB} {-1 SUB}]
}
```

You can also make use of built-in Lexer functions: `ReadInt`, `ReadIntOrFloat`, `ReadBetween`
```go
func main3() {
	type Token struct {
		Type  string
		Value string
	}

	const (
		TOKEN_ADD = iota
		TOKEN_INT
		TOKEN_FLOAT
	)

	lex := lexer.NewLexer(lexer.Options{
		Tokens: []int{TOKEN_ADD, TOKEN_FLOAT, TOKEN_INT},
		Matchers: map[int]lexer.MatcherFunc{
			TOKEN_ADD: func(l *lexer.Lexer) (lexer.Token, bool) {
				ok := l.ReadChar('+')
				if !ok {
					return Token{}, false
				}
				return Token{"TOKEN_ADD", "+"}, true
			},
			TOKEN_FLOAT: func(l *lexer.Lexer) (lexer.Token, bool) {
				str, ok := l.ReadIntOrFloat()
				if !ok {
					return Token{}, false
				}
				return Token{"TOKEN_FLOAT", str}, true
			},
		},
	})

	lexerTokens, _ := lex.Lex("5+5.6")

	fmt.Println(lexerTokens)
	// [{TOKEN_FLOAT 5} {TOKEN_ADD +} {TOKEN_FLOAT 5.6}]
}
```