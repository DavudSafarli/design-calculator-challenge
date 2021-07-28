package calculator

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		input string
		want  []Token
	}{
		{"123456789", []Token{{NUM, "123456789"}}},
		{"  1  ", []Token{{NUM, "1"}}},
		{"+", []Token{{ADD, ""}}},
		{"-", []Token{{SUB, ""}}},
		{"*", []Token{{MUL, ""}}},
		{"/", []Token{{DIV, ""}}},
		{"^", []Token{{POW, ""}}},
		{"()", []Token{{L_PAR, ""}, {R_PAR, ""}}},
		{"1+2", []Token{{NUM, "1"}, {ADD, ""}, {NUM, "2"}}},
		{"1-2", []Token{{NUM, "1"}, {SUB, ""}, {NUM, "2"}}},
		{
			"  1  +  2  +  3  ",
			[]Token{{NUM, "1"}, {ADD, ""}, {NUM, "2"}, {ADD, ""}, {NUM, "3"}},
		},
		{
			"7-7/7+7*14",
			[]Token{{NUM, "7"}, {SUB, ""}, {NUM, "7"}, {DIV, ""}, {NUM, "7"}, {ADD, ""}, {NUM, "7"}, {MUL, ""}, {NUM, "14"}},
		},
	}
	for _, tt := range tests {
		testName := fmt.Sprint("Tokenizing", tt.input)
		t.Run(testName, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			tokenizer := Tokenizer{}

			got, err := tokenizer.Tokenize(reader)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}
