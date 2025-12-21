package main

import (
	"errors"
	"math"
	"reflect"
)

type Callable interface {
	LispExp
	Call(args []any) (any, error)
}

type Procedure struct {
	params []string
	body   *ConsCell
	env    *Env
}

func (p *Procedure) isLispExp() {}

func (p *Procedure) Call(args []any) (any, error) {
	env := NewEnv(p.params, args, p.env)
	return Eval(p.body, env)
}

type BuiltIn struct {
	body func(args []any) (any, error)
}

func (b *BuiltIn) isLispExp() {}

func (b *BuiltIn) Call(args []any) (any, error) {
	return b.body(args)
}

func equalOp(args []any) (any, error) {
	if len(args) != 2 {
		return nil, errors.New("expected to arguments to equal?")
	}
	if reflect.DeepEqual(args[0], args[1]) {
		return true, nil
	}
	return false, nil
}

func addOp(args []any) (any, error) {
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
	return total, nil
}
func multOp(args []any) (any, error) {
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
	return product, nil
}
func divOp(args []any) (any, error) {
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
	return result, nil
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
