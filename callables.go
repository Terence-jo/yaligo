package main

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

type Callable interface {
	LispExp
	Call(args *ConsCell) (LispExp, error)
}

type Procedure struct {
	params *ConsCell
	body   LispExp
	env    *Env
}

func (p *Procedure) isLispExp() {}

func (p *Procedure) String() string {
	return fmt.Sprintf("(lambda (%s) (%s))", p.params, p.body)
}

func (p *Procedure) Call(args *ConsCell) (LispExp, error) {
	env, err := NewEnv(p.params, args, p.env)
	if err != nil {
		return nil, err
	}
	return Eval(p.body, env)
}

type BuiltIn struct {
	body func(args *ConsCell) (LispExp, error)
}

func (b *BuiltIn) isLispExp() {}

func (b *BuiltIn) String() string {
	return "yaligo built-in"
}

func (b *BuiltIn) Call(args *ConsCell) (LispExp, error) {
	return b.body(args)
}

func equalOp(args *ConsCell) (LispExp, error) {
	if args.Cdr == nil || args.Cdr.Cdr != nil {
		return nil, errors.New("expected two arguments to equal?")
	}
	if reflect.DeepEqual(args.Car, args.Cdr.Car) {
		return &SymbolAtom{TRUE}, nil
	}
	return &SymbolAtom{FALSE}, nil
}

func addOp(args *ConsCell) (LispExp, error) {
	if args.Car == nil {
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
	return &NumberAtom{total}, nil
}

func multOp(args *ConsCell) (LispExp, error) {
	if args.Car == nil {
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
	return &NumberAtom{product}, nil
}
func divOp(args *ConsCell) (LispExp, error) {
	if args.Car == nil {
		return nil, errors.New("expected arguments to *")
	}
	// scheme behaviour is to return the inverse if given only one arg
	var accumulator float64
	div := func(x, y float64) float64 {
		return x / y
	}
	if args.Cdr == nil {
		accumulator = 1.0
		inverse, err := reduceNums(args, div, accumulator)
		if err != nil {
			return nil, err
		}
		return &NumberAtom{Data: inverse}, nil
	}
	accumulator, err := numToFloat(args.Car)
	if err != nil {
		return nil, err
	}
	result, err := reduceNums(args.Cdr, div, accumulator)
	if err != nil {
		return nil, err
	}
	return &NumberAtom{result}, nil
}

func reduceNums(nums *ConsCell, fn func(float64, float64) float64, accumulator float64) (float64, error) {
	for nums != nil {
		num, ok := nums.Car.(*NumberAtom)
		if !ok {
			return math.NaN(), errors.New("expected all arguments to be numbers")
		}
		floatVal, err := numToFloat(num.Data)
		if err != nil {
			return math.NaN(), err
		}
		accumulator = fn(accumulator, floatVal)
		nums = nums.Cdr
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
