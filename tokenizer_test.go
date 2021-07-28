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
		{
			"1",
			[]Token{{NUM, "1"}},
		},
		{
			"  1  ",
			[]Token{{NUM, "1"}},
		},
		{
			"+",
			[]Token{{ADD, ""}},
		},
		{
			"1+2",
			[]Token{{NUM, "1"}, {ADD, ""}, {NUM, "2"}},
		},
		{
			"  1  +  2  +  3  ",
			[]Token{{NUM, "1"}, {ADD, ""}, {NUM, "2"}, {ADD, ""}, {NUM, "3"}},
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
