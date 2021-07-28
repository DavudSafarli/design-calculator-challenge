package calculator

type Calculatable interface {
	Calculate() float64
}

type NumNode struct {
	Value float64
}

func (n NumNode) Calculate() float64 {
	return n.Value
}

type AddNode struct {
	Left  Calculatable
	Right Calculatable
}

func (n AddNode) Calculate() float64 {
	return n.Left.Calculate() + n.Right.Calculate()
}
