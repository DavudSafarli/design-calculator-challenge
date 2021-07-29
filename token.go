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
	return t.Type == NUM && t.Value != ""
}

func (t Token) IsOP() bool {
	return t.Type == ADD || t.Type == SUB || t.Type == MUL || t.Type == DIV || t.Type == POW
}

func (t Token) IsAddOP() bool {
	return t.Type == ADD
}
func (t Token) IsSubOP() bool {
	return t.Type == SUB
}
func (t Token) IsMulOP() bool {
	return t.Type == MUL
}
func (t Token) IsDivOP() bool {
	return t.Type == DIV
}
func (t Token) IsPowOP() bool {
	return t.Type == POW
}
func (t Token) IsParacentesis() bool {
	return t.Type == L_PAR || t.Type == R_PAR
}
func (t Token) IsLeftParacentesis() bool {
	return t.Type == L_PAR
}
func (t Token) IsRightParacentesis() bool {
	return t.Type == R_PAR
}
