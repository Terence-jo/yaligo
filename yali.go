package main

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

func addOp(x float64, y float64) float64 {
	return x + y
}
func multOp(x float64, y float64) float64 {
	return x * y
}
func divOp(x float64, y float64) float64 {
	return x / y
}

// Eval() will evaluate the list it is passed. It will return a number
// at this stage, but in the future will need to be flexible enough to
// return a number, string, or list. It will just need to return
// any I think. It may have evaluation errors.
//
// Post-calculator considerations:
// Where will the environment come from? If I model the environment as
// a simple map, what are the drawbacks? I think the main one is that
// it could be hard to extricate the values from an inner function from the environment
// for an outer function. There are ways around this though, the first to come
// to mind is to model the environment as a map of maps, but perhaps a tree
// of maps would be more appropriate. That would allow clearer definition
// right?
//
// See about structuring the data types to be able to approach this with actual lisp-like
// semantics. That would be nice.
func Eval(item Linker) (any, error) {
	// I had forgotten that type switches could be so elegant in Go. this is nice!
	switch item := item.(type) {
	case *IntItem:
		floatRet := float64(item.Data)
		return floatRet, nil
	case *FloatItem:
		return item.Data, nil
	case *SymbolItem:
		op := item.Data
		switch op {
		case "+":
			return reduceNums(item, addOp, 0.0)
		case "*":
			return reduceNums(item, multOp, 1.0)
		case "/":
			item := item.Next()
			left, err := Eval(item)
			if err != nil {
				return nil, err
			}
			floatLeft, ok := left.(float64)
			if !ok {
				return math.NaN(), errors.New("expected number as operand to /")
			}
			return reduceNums(item.Next(), divOp, floatLeft)
		}
	case *ListItem:
		return Eval(item.Data.Car())
	}
	return nil, errors.New("type did not match")
}

func reduceNums(item Linker, fn func(float64, float64) float64, initVal float64) (float64, error) {
	total := initVal
	for item.Next() != nil {
		right, err := Eval(item.Next())
		if err != nil {
			return math.NaN(), err
		}
		floatRight, ok := right.(float64)
		if !ok {
			return math.NaN(), errors.New("expected number as operand")
		}
		total = fn(total, floatRight)
		item = item.Next()
	}
	return total, nil
}

func readFromTokens(tokens []Token, pos int) (Linker, int, error) {
	if len(tokens) == 0 {
		return nil, 0, errors.New("unexpected EOF")
	}
	token := tokens[pos]
	switch token.Class {
	case OPEN:
		// start a list and recurse to fill it
		pos++
		list := NewList()
		cur := list.head
		for tokens[pos].Class != CLOSE {
			next, new_pos, err := readFromTokens(tokens, pos)
			if err != nil {
				return nil, 0, err
			}
			cur.SetNext(next)
			cur = next
			pos = new_pos
		}
		expr := ListItem{Data: list}
		return &expr, pos, nil
	case CLOSE:
		return nil, 0, errors.New("unexpected )")
	}
	atom, err := atom(token)
	if err != nil {
		return nil, 0, err
	}
	pos++
	return atom, pos, nil
}

func tokenise(chars string) []string {
	chars = strings.Replace(
		chars, "(", " ( ", -1,
	)
	chars = strings.Replace(
		chars, ")", " ) ", -1,
	)
	tokens := []string{}
	for _, token := range strings.Split(chars, " ") {
		if len(token) > 0 {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func atom(token Token) (Linker, error) {
	if token.Class == NUMBER {
		// is this acceptable as a way to try parsing the numbers? and do I
		// just want to make everything float?
		intval, err := strconv.ParseInt(token.Lit, 10, 64)
		if err != nil {
			floatval, err := strconv.ParseFloat(token.Lit, 64)
			if err != nil {
				return nil, err
			}
			return &FloatItem{Data: floatval}, nil
		}
		return &IntItem{Data: intval}, nil
	}
	return &SymbolItem{Data: token.Lit}, nil
}
