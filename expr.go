package main

// An Expr is an arithmetic expression.
type Expr interface {
	// Eval returns the value of this Expr in the environment env.
	Eval() (float64, error)
	// Expr is a Stringer too
	String() string
	// Len returns the number of symbols of the expression. (A number is just one symbol.)
	Len() int
}
