package main

import (
	"fmt"
)

// A num is a floating number
type num float64

func (f num) Eval() (float64, error) {
	return float64(f), nil
}
func (f num) String() string {
	return fmt.Sprintf("%.2f", f)
}
func (f num) Len() int {
	return 1
}

// A unary is an operator with only one operand
type unary struct {
	op rune // one of '+', '-'
	x  Expr
}

func (u unary) String() string {
	return fmt.Sprintf("%s%s", string(u.op), u.x)
}

func (u unary) Eval() (float64, error) {
	x, err := u.x.Eval()
	if err != nil {
		return 0, fmt.Errorf("evaluation of operand x = %v in unary failed: %s", u.x, err)
	}
	switch u.op {
	case '+':
		return +x, nil
	case '-':
		return -x, nil
	}
	return 0, fmt.Errorf("unsupported unary operator: %q", u.op)
}

func (u unary) Len() int {
	return u.x.Len() + 1
}

// A binary is an operator with two operands
type binary struct {
	op   rune // one of '+', '-', '*', '/'
	x, y Expr
}

func (u binary) String() string {
	return fmt.Sprintf("%s %s %s", u.x, string(u.op), u.y)
}

func (b binary) Eval() (float64, error) {
	x, err := b.x.Eval()
	if err != nil {
		return 0, fmt.Errorf("evaluation of operand x = %v in binary failed: %s", b.x, err)
	}
	y, err := b.y.Eval()
	if err != nil {
		return 0, fmt.Errorf("evaluation of operand y = %v in binary failed: %s", b.y, err)
	}
	switch b.op {
	case '+':
		return x + y, nil
	case '-':
		return x - y, nil
	case '*':
		return x * y, nil
	case '/':
		if y == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return x / y, nil
	default:
		return 0, fmt.Errorf("unsupported binary operator: %q", b.op)
	}
}

func (b binary) Len() int {
	return b.x.Len() + b.y.Len() + 1
}
