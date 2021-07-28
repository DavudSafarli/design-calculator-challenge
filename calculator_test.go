package calculator

import "testing"

func TestCalculator(t *testing.T) {
	c := NewCalculator()
	actual := c.Eval("1+2")
	var expected float64 = 3

	if actual != expected {
		t.Fatalf("expected: %v\nactual: %v", expected, actual)
	}

}
