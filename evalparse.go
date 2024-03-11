package main

// The functions in this package can be used to parse an arthimetic expression in text form into an ast-like data structure.
// We use a lexer and a parse algorithm for symbolic expressions; both ideas come from Donovan & Kernighan (2016)

import (
	"fmt"
	"io"
	"strconv"
	"text/scanner"
)

// EvalParse parses the content from the input reader as an arithmetic expression.
// It uses an adaptation of a parse algorithm for symbolic expressions by D&K(2016)
// In addition to its counterpart Parse(), it makes evaluation in place of parsed operands.
// This way, the returned Expr is in fact a num.
func EvalParse(r io.Reader) (Expr, error) {
	lex := new(lexer)
	lex.scan.Init(r)
	lex.scan.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats
	lex.next() // initial lookahead
	e, err := evalparseExpr(lex)
	if err != nil {
		return nil, fmt.Errorf("could not parse %s: %s", lex, err)
	}
	if lex.token != scanner.EOF {
		return nil, fmt.Errorf("unexpected %s", lex)
	}

	return e, nil
}

func evalparseExpr(lex *lexer) (Expr, error) { return evalparseBinary(lex, 1) }

// evalparseBinary stops when it encounters an
// operator of lower prio than prio0.
func evalparseBinary(lex *lexer, prio0 int) (Expr, error) {
	left, err := evalparseUnary(lex)
	if err != nil {
		return nil, fmt.Errorf("could not parse expression in unary %s: %s", lex, err)
	}
	for prio := priority(lex.token); prio >= prio0; prio-- {
		for priority(lex.token) == prio {
			op := lex.token
			lex.next() // consume operator
			right, err := evalparseBinary(lex, prio+1)
			if err != nil {
				return nil, fmt.Errorf("could not parse expression in unary %s: %s", lex, err)
			}
			leftEval, _ := left.Eval()
			left = binary{op, num(leftEval), right}
			// left = binary{op, left, right}
		}
	}
	leftEval, _ := left.Eval()
	return num(leftEval), nil
}

func evalparseUnary(lex *lexer) (Expr, error) {
	if lex.token == '+' || lex.token == '-' {
		op := lex.token
		lex.next() // consume '+' or '-'
		e, err := evalparseUnary(lex)
		if err != nil {
			return nil, fmt.Errorf("could not parse expression in unary %s: %s", lex, err)
		}
		eEval, _ := e.Eval()
		return unary{op, num(eEval)}, nil
		// return unary{op, e}, nil
	}
	return evalparsePrimary(lex)
}

func evalparsePrimary(lex *lexer) (Expr, error) {
	switch lex.token {
	case scanner.Int, scanner.Float:
		f, err := strconv.ParseFloat(lex.text(), 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse the float number %s: %s", lex, err)
		}
		lex.next() // consume number
		return num(f), nil

	case '(':
		lex.next() // consume '('
		e, err := evalparseExpr(lex)
		eEval, _ := e.Eval()
		if err != nil {
			return nil, fmt.Errorf("could not parse the symbol %s: %s", lex, err)
		}
		if lex.token != ')' {
			return nil, fmt.Errorf("got %s, want ')'", lex)
		}
		lex.next() // consume ')'
		return num(eEval), nil
	}
	return nil, fmt.Errorf("unexpected %s", lex)
}
