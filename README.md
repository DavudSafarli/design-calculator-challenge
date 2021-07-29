# Calculator

Simple calculator application, written in Go.
Supports basic `BODMAS` operations.


Example usage:
```
func main() {
    c := calculator.New()
    res, err := c.Eval("5+5")
    fmt.Println(res, err)
}
```

## How it works

It goes through multiple steps to calculate the input expression.

1. First step is `lexical analysis` of the given input, and converting it into a slice of Token. You can think of a Token as a `keyword`, or a character that has a meaning. It can also hold any value related to itself. In the context of this application, some of the Tokens are :
   - NUM
   - ADD
   - SUB
   - MUL
   - DIV
   - POW
   - L_PAR
   - R_PAR

    As an example, expression "1+2" after lexical analysis produces:
    `[Token(NUM, 1), Token(ADD), Token(NUM, 2)]`

2. Second step is validating the tokenized input find and report errors and error positions. For example:
`"(*4" -> [LEFT_PAR, MUL, NUM]` is an invalid expression. The position(index) where it starts being invalid is 1(MUL), and ends at 2, because we cannot have Multiplication after the opening-paracantesis. Here, we report the Error(`ErrOperationAfterLeftParacantesis`), the start and endposition, so that we can show the incorrect part to the end-user.


3. Parsing the Tokens, and building the [Expression Tree](https://www.geeksforgeeks.org/expression-tree/). The way we build the expression tree is parsing the Tokens from `Infix notation` to `Postfix` and then, `Postfix` to `Expression Tree`. Here in this application, this steps happen at the same time. You can learn more about `Infix to Postfix` from [this video](https://youtu.be/PAceaOSnxQs). You can read this [geeksforgeeks article on Program to convert Infix notation to Expression Tree](https://www.geeksforgeeks.org/program-to-convert-infix-notation-to-expression-tree/) to get learn more about the algorithm.

4. Evaluating the Expression tree `recursively` and getting the result.