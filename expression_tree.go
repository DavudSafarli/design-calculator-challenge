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

type SubNode struct {
	Left  Calculatable
	Right Calculatable
}

func (n SubNode) Calculate() float64 {
	return n.Left.Calculate() - n.Right.Calculate()
}

type MulNode struct {
	Left  Calculatable
	Right Calculatable
}

func (n MulNode) Calculate() float64 {
	return n.Left.Calculate() * n.Right.Calculate()
}

type DivNode struct {
	Left  Calculatable
	Right Calculatable
}

func (n DivNode) Calculate() float64 {
	return n.Left.Calculate() / n.Right.Calculate()
}
