package calculator_test

import (
	"fmt"
	"testing"

	calculator "github.com/DavudSafarli/design-calculator-challenge"
	"github.com/DavudSafarli/design-calculator-challenge/lexer"
)

func TestValidExpressions(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"1+2", 3},
		{"1-2", -1},
		{"2+3-3", 2},
		{"3+2-4-2-2", -3},
		{"9-3+2-4-2-2", 0},
		{"1*2*9", 18},
		{"5* 2-3", 7},
		{"5+2*3", 11},
		{"5/2-3", -0.5},
		{"36/6*3-4+8", 22},
		{"20+12/4*3-2*2", 25},
		{"7-7/7+7*14", 104},
		{"2*3+4", 10},
		{"2*(3+4)", 14},
		{"7-7/7+7*(7+7)", 104},
		{"2^3", 8},
		{"2^3+2", 10},
		{"2^3+2", 10},
		{"25-48/6+12*2", 41},
		{"3*(4*5^2)/6+7-8", 49},
		{"(0-10)/(20/2^2*5/5)*8-2", -18},
		{"2.5*4.0+2", 12},
		{"2.5(4.0+2)", 15},
		{"2.5((4.0+2))", 15},
		{"2.5+((4.0+2))", 8.5},
	}

	for _, tt := range tests {
		testName := fmt.Sprint("Calculating ", tt.input)
		t.Run(testName, func(t *testing.T) {
			calc := calculator.New()
			actual, evalErr := calc.Eval(tt.input)

			if evalErr != nil {
				t.Fatalf("\nexpected: nil\nactual  : %v", evalErr)
			}

			if actual != tt.want {
				t.Fatalf("\nexpected: %v\nactual  : %v", tt.want, actual)
			}
		})
	}

}

func TestInvalidExpressions(t *testing.T) {
	tests := []struct {
		input string
		err   error
	}{
		{"1+)", calculator.EvalError{calculator.ErrOperationBeforeRightParacentesis, 2, 3}},
		{"(3))", calculator.EvalError{calculator.ErrInconsistentParacentesisCount, 3, 4}},
		{"(5+/4", calculator.EvalError{calculator.Err2Operators, 3, 4}},
		{"((*5+4", calculator.EvalError{calculator.ErrOperationAfterLeftParacantesis, 2, 3}},
		{"*5+4", calculator.EvalError{calculator.ErrCannotStartWithOperator, 0, 1}},
		{"(5", calculator.EvalError{calculator.ErrInconsistentParacentesisCount, -1, -1}},
		{"  1.2+*", calculator.EvalError{calculator.Err2Operators, 6, 7}},
	}

	for _, tt := range tests {
		testName := fmt.Sprint("Calculating ", tt.input)
		t.Run(testName, func(t *testing.T) {
			c := calculator.New()
			_, evalErr := c.Eval(tt.input)
			if evalErr != tt.err {
				t.Fatalf("\nexpected: %v\nactual  : %v", tt.err, evalErr)
			}
		})
	}
}

func TestInvalidInputs(t *testing.T) {
	tests := []struct {
		input string
		err   error
	}{
		{"5+&&??", lexer.UnknownSymbolError{Symbol: '&'}},
	}

	for _, tt := range tests {
		testName := fmt.Sprint("Calculating ", tt.input)
		t.Run(testName, func(t *testing.T) {
			c := calculator.New()
			_, evalErr := c.Eval(tt.input)
			if evalErr != tt.err {
				t.Fatalf("\nexpected: %v\nactual  : %v", tt.err, evalErr)
			}
		})
	}
}
