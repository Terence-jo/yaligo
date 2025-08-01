package main

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

type Env map[string]any

func standardEnv() Env {
	return Env{
		"+": addOp,
		"*": multOp,
		"/": divOp,
	}
}

var globalEnv = standardEnv()

func addOp(args []any) (LispExp, error) {
	if len(args) == 0 {
		return nil, errors.New("expected arguments to +")
	}
	accumulator := 0.0
	add := func(x float64, y float64) float64 {
		return x + y
	}
	total, err := reduceNums(args, add, accumulator)
	if err != nil {
		return nil, err
	}
	ret := &FloatAtom{data: total}
	return ret, nil
}
func multOp(args []any) (LispExp, error) {
	if len(args) == 0 {
		return nil, errors.New("expected arguments to *")
	}
	accumulator := 1.0
	mult := func(x float64, y float64) float64 {
		return x * y
	}
	product, err := reduceNums(args, mult, accumulator)
	if err != nil {
		return nil, err
	}
	ret := &FloatAtom{data: product}
	return ret, nil
}
func divOp(args []any) (LispExp, error) {
	if len(args) == 0 {
		return nil, errors.New("expected arguments to *")
	}
	// scheme behaviour is to return the inverse if given only one arg
	var accumulator float64
	div := func(x, y float64) float64 {
		return x / y
	}
	if len(args) == 1 {
		accumulator = 1.0
		inverse, err := reduceNums(args, div, accumulator)
		if err != nil {
			return nil, err
		}
		return &FloatAtom{data: inverse}, nil
	}
	//
	accumulator, err := numToFloat(args[0])
	if err != nil {
		return nil, err
	}
	result, err := reduceNums(args[1:], div, accumulator)
	if err != nil {
		return nil, err
	}
	ret := &FloatAtom{data: result}
	return ret, nil
}

func reduceNums(nums []any, fn func(float64, float64) float64, accumulator float64) (float64, error) {
	for _, num := range nums {
		floatVal, err := numToFloat(num)
		if err != nil {
			return math.NaN(), err
		}
		accumulator = fn(accumulator, floatVal)
	}
	return accumulator, nil
}

func numToFloat(num any) (float64, error) {
	switch val := num.(type) {
	case int64:
		return float64(val), nil
	case float64:
		return val, nil
	default:
		return math.NaN(), errors.New("expected number as operand")
	}
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
func Eval(exp LispExp, env Env) (any, error) {
	// I had forgotten that type switches could be so elegant in Go. this is nice!
	switch exp := exp.(type) {
	case *IntAtom:
		intVal := exp.Value().(int64)
		floatRet := float64(intVal)
		return floatRet, nil
	case *FloatAtom:
		floatRet := exp.Value().(float64)
		return floatRet, nil
	case *SymbolAtom:
		// this case should just be for evaluating a symbol in the environment
		symbol := exp.Value().(string)
		return env[symbol], nil
	case *ListExp:
		// assuming this means it is a procedure call:
		car, err := Eval(exp.Car(), env)
		if err != nil {
			return nil, err
		}
		proc, ok := car.(func([]any) (LispExp, error))
		if !ok {
			return nil, errors.New("expected procedure name at head of list")
		}
		argExps := exp.Cdr()
		var argVals []any
		for {
			arg, err := Eval(argExps.Car(), env)
			if err != nil {
				return nil, err
			}
			argVals = append(argVals, arg)
			if argExps.Cdr() == nil {
				break
			}
			argExps = argExps.Cdr()
		}
		ret, err := proc(argVals)
		if err != nil {
			return nil, err
		}
		// Returning a LispExp from procedures and evaluating the return avoids another
		// type switch, but it does increase recursion depth momentarily...
		val, err := Eval(ret, env)
		if err != nil {
			return nil, err
		}
		return val, nil
	}
	return nil, errors.New("type did not match")
}

func readFromTokens(tokens []Token, pos int) (LispExp, int, error) {
	if len(tokens) == 0 {
		return nil, 0, errors.New("unexpected EOF")
	}
	token := tokens[pos]
	switch token.Class {
	case OPEN:
		// start a list and recurse to fill it
		pos++
		list := NewList()
		var cur LispExp = list
		for tokens[pos].Class != CLOSE {
			next, new_pos, err := readFromTokens(tokens, pos)
			if err != nil {
				return nil, 0, err
			}
			cur.SetNext(next)
			cur = next
			pos = new_pos
		}
		return list, pos, nil
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
	chars = strings.ReplaceAll(
		chars, "(", " ( ",
	)
	chars = strings.ReplaceAll(
		chars, ")", " ) ",
	)
	tokens := []string{}
	for _, token := range strings.Split(chars, " ") {
		if len(token) > 0 {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func atom(token Token) (LispExp, error) {
	if token.Class == NUMBER {
		// is this acceptable as a way to try parsing the numbers? and do I
		// just want to make everything float?
		intval, err := strconv.ParseInt(token.Lit, 10, 64)
		if err != nil {
			floatval, err := strconv.ParseFloat(token.Lit, 64)
			if err != nil {
				return nil, err
			}
			return &FloatAtom{data: floatval}, nil
		}
		return &IntAtom{data: intval}, nil
	}
	return &SymbolAtom{data: token.Lit}, nil
}
