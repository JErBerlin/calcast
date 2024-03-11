package main

// The functions in this package can be used to parse an arthimetic expression in text form into an ast-like data structure.
// We use a lexer and a parse algorithm for symbolic expressions; both ideas come from Donovan & Kernighan (2016)

import (
	"fmt"
	"io"
	"strconv"
	"text/scanner"
)

// lexer
type lexer struct {
	scan  scanner.Scanner
	token rune // current token, used as lookahead
}

func (lex *lexer) next()        { lex.token = lex.scan.Scan() } // consumes and stores token
func (lex *lexer) text() string { return lex.scan.TokenText() } // return last scanned token as text

// String returns a string describing the current state of the lexer (the current token)
// for use in errors.
func (lex *lexer) String() string {
	switch lex.token {
	case scanner.EOF:
		return "end of file"
	case scanner.Ident:
		return fmt.Sprintf("identifier %s", lex.text())
	case scanner.Int, scanner.Float:
		return fmt.Sprintf("number %s", lex.text())
	}
	return fmt.Sprintf("%q", rune(lex.token)) // any other rune
}

func priority(op rune) int {
	switch op {
	case '*', '/':
		return 2
	case '+', '-':
		return 1
	}
	return 0
}

// Parse parses the content from the input reader as an arithmetic expression.
// It uses lazy loading. The buffering management is done by the scanner in the lexer.
func Parse(r io.Reader) (Expr, error) {
	lex := new(lexer)
	lex.scan.Init(r)

	// configure the lexer
	// recognise as tokens: symbols, integers and floats
	lex.scan.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats

	lex.next() // initial lookahead
	e, err := parseExpr(lex)
	if err != nil {
		return nil, fmt.Errorf("could not parse %s: %s", lex, err)
	}
	if lex.token != scanner.EOF {
		return nil, fmt.Errorf("unexpected %s", lex)
	}

	return e, nil
}

// parseExpr is just an entry point to parseBinary with a low operator priority of 1 
// this represents a sum A + B, or a rest A - B
func parseExpr(lex *lexer) (Expr, error) { return parseBinary(lex, 1) }

// parseBinary parses a binary operation with its operands: -A + (B) or -A * (B)
// it stops when it encounters an operator of lower prio than prio0
func parseBinary(lex *lexer, prio0 int) (Expr, error) {
	left, err := parseUnary(lex)
	if err != nil {
		return nil, fmt.Errorf("could not parse expression in unary %s: %s", lex, err)
	}

	for prio := priority(lex.token); prio >= prio0; prio-- { 
		for priority(lex.token) == prio { 
			op := lex.token
			lex.next() // consume operator and look ahead
			right, err := parseBinary(lex, prio+1)
			if err != nil {
				return nil, fmt.Errorf("could not parse expression in unary %s: %s", lex, err)
			}
			left = binary{op, left, right}
		}
	}
	return left, nil
}

// parses a signed number or a signed parenthesis: -A or -(...)
func parseUnary(lex *lexer) (Expr, error) {
	if lex.token == '+' || lex.token == '-' {
		op := lex.token
		lex.next() // consume '+' or '-'
		e, err := parseUnary(lex)
		if err != nil {
			return nil, fmt.Errorf("could not parse expression in unary %s: %s", lex, err)
		}
		return unary{op, e}, nil
	}
	// parse number or parenthesis group after the sign
	return parsePrimary(lex)
}

// parsePrimary parses a number or a parenthesis group: N or (...)
func parsePrimary(lex *lexer) (Expr, error) {
	switch lex.token {

	// parse an integer or a float number
	case scanner.Int, scanner.Float:
		f, err := strconv.ParseFloat(lex.text(), 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse the float number %s: %s", lex, err)
		}
		lex.next() // consume number
		return num(f), nil

	case '(':
		lex.next() // consume '('
		
		// parse expression inside parenthesis
		e, err := parseExpr(lex)
		if err != nil {
			return nil, fmt.Errorf("could not parse the symbol %s: %s", lex, err)
		}
		
		if lex.token != ')' {
			return nil, fmt.Errorf("got %s, want ')'", lex)
		}
		lex.next() // consume ')'
		
		return e, nil
	}
	return nil, fmt.Errorf("unexpected %s", lex)
}
