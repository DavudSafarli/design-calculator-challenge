package calculator

// TOKENS
const (
	NUM = iota
	ADD
	SUB
	MUL
	DIV
	POW
	L_PAR
	R_PAR
)

type Token struct {
	Type  int
	Value string
}

func (t Token) IsNum() bool {
	return t.Type == NUM
}

func (t Token) IsOperand() bool {
	return t.Type == ADD || t.Type == SUB || t.Type == MUL || t.Type == DIV || t.Type == POW
}

func (t Token) IsAddOperand() bool {
	return t.Type == ADD
}
func (t Token) IsSubOperand() bool {
	return t.Type == SUB
}
func (t Token) IsMulOperand() bool {
	return t.Type == MUL
}
func (t Token) IsDivOperand() bool {
	return t.Type == DIV
}
func (t Token) IsPowOperand() bool {
	return t.Type == POW
}
func (t Token) IsLeftParacentesis() bool {
	return t.Type == L_PAR
}
func (t Token) IsRightParacentesis() bool {
	return t.Type == R_PAR
}
