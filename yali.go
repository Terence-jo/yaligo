package main

import (
	"errors"
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
	total := 0.0
	for _, arg := range args {
		switch val := arg.(type) {
		case int64:
			total += float64(val)
		case float64:
			total += val
		default:
			return nil, errors.New("expected number as arg to +")
		}
	}
	ret := &FloatAtom{data: total}
	return ret, nil
}
func multOp(args []any) (LispExp, error) {
	if len(args) == 0 {
		return nil, errors.New("expected arguments to *")
	}
	product := 1.0
	for _, arg := range args {
		switch val := arg.(type) {
		case int64:
			product *= float64(val)
		case float64:
			product *= val
		default:
			return nil, errors.New("expected number as arg to *")
		}
	}
	ret := &FloatAtom{data: product}
	return ret, nil
}
func divOp(args []any) (LispExp, error) {
	if len(args) == 0 {
		return nil, errors.New("expected arguments to *")
	}
	// why did I change the language here, using atom?
	atom := args[0]
	quotient, ok := atom.(Atom).Value().(float64)
	if !ok {
		return nil, errors.New("expected number atom as argument to /")
	}
	// scheme behaviour is to reciprocate the value if given only one argument
	if len(args) == 1 {
		if quotient == 0.0 {
			return nil, errors.New("expected non-zero number atom as single argument to /")
		}
		return &FloatAtom{data: 1 / quotient}, nil
	}
	for _, arg := range args[1:] {
		switch val := arg.(type) {
		case int64:
			quotient /= float64(val)
		case float64:
			quotient /= val
		default:
			return nil, errors.New("expected number as arg to *")
		}
	}
	ret := &FloatAtom{data: quotient}
	return ret, nil
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

// func reduceNums(item LispExp, fn func(float64, float64) float64, initVal float64) (float64, error) {
// 	total := initVal
// 	for item.Next() != nil {
// 		right, err := Eval(item.Next())
// 		if err != nil {
// 			return math.NaN(), err
// 		}
// 		floatRight, ok := right.(float64)
// 		if !ok {
// 			return math.NaN(), errors.New("expected number as operand")
// 		}
// 		total = fn(total, floatRight)
// 		item = item.Next()
// 	}
// 	return total, nil
// }

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
