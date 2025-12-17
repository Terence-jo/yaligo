package main

import (
	"errors"
	"math"
	"reflect"
	"strconv"
	"strings"
)

var globalEnv = standardEnv()

func equalOp(args []any) (LispExp, error) {
	if len(args) != 2 {
		return nil, errors.New("expected to arguments to equal?")
	}
	if reflect.DeepEqual(args[0], args[1]) {
		return &SymbolAtom{Data: "#t"}, nil
	}
	return &SymbolAtom{Data: "#f"}, nil
}

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
	ret := &FloatAtom{Data: total}
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
	ret := &FloatAtom{Data: product}
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
		return &FloatAtom{Data: inverse}, nil
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
	ret := &FloatAtom{Data: result}
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
func Eval(exp LispExp, env *Env) (any, error) {
	// I had forgotten that type switches could be so elegant in Go. this is nice!
	switch exp := exp.(type) {
	case *IntAtom:
		intVal := exp.Data
		floatRet := float64(intVal)
		return floatRet, nil
	case *FloatAtom:
		floatRet := exp.Data
		return floatRet, nil
	case *SymbolAtom:
		// this case should just be for evaluating a symbol in the environment
		symbol := exp.Data
		return env.Find(symbol), nil
	case *ConsCell:
		car := exp.Car
		exp = exp.Cdr
		// unquoted list, this means it is a syntactic form or procedure call:
		symbol, ok := car.(*SymbolAtom)
		if !ok {
			return nil, errors.New("expected symbol at head of unquoted list")
		}
		switch symbol.Data {
		case "if":
			// chucking this here to test as a helper.
			test := exp.Car
			conseq := exp.Cdr.Car
			alt := exp.Cdr.Cdr.Car
			exp := exp.Cdr.Cdr.Cdr
			if exp != nil {
				return nil, errors.New("too many arguments to 'if', expected test, conseq, alt")
			}
			testRes, err := Eval(test, env)
			if err != nil {
				return nil, err
			}
			resBool, ok := testRes.(bool)
			if !ok {
				return nil, errors.New("expected test of 'if' to evaluate to a boolean")
			}
			if resBool {
				return Eval(conseq, env)
			}
			return Eval(alt, env)
		case "cond":
		case "define":
		default:
			procExp, err := Eval(symbol, env)
			if err != nil {
				return nil, err
			}
			proc, ok := procExp.(func([]any) (LispExp, error))
			if !ok {
				return nil, errors.New("expected procedure name at head of list")
			}
			var argVals []any
			for {
				arg, err := Eval(exp.Car, env)
				if err != nil {
					return nil, err
				}
				argVals = append(argVals, arg)
				if exp.Cdr == nil {
					break
				}
				exp = exp.Cdr
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
		pos++
		if len(tokens) <= pos {
			return nil, 0, errors.New("unclosed list: unexpected EOF")
		}
		// Handle empty list '()'
		if tokens[pos].Class == CLOSE {
			return nil, pos + 1, nil
		}

		// Read the first element to create the head of the list.
		car, newPos, err := readFromTokens(tokens, pos)
		if err != nil {
			return nil, 0, err
		}
		head := Cons(car, nil)
		tail := head
		pos = newPos
		// start a list and recurse to fill it
		for tokens[pos].Class != CLOSE {
			car, new_pos, err := readFromTokens(tokens, pos)
			if err != nil {
				return nil, 0, err
			}
			newCell := Cons(car, nil)
			tail.Cdr = newCell
			tail = newCell
			pos = new_pos
		}
		return head, pos, nil
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
			return &FloatAtom{Data: floatval}, nil
		}
		return &IntAtom{Data: intval}, nil
	}
	return &SymbolAtom{Data: token.Lit}, nil
}
