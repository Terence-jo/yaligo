package main

import "errors"

// symbols in standard environment
const (
	TRUE  string = "#t"
	FALSE string = "#f"
	EQUAL string = "eq?"
	ADD   string = "+"
	MULT  string = "*"
	DIV   string = "/"
)

type Env struct {
	local map[string]LispExp
	outer *Env
}

func NewEnv(params *ConsCell, args *ConsCell, outer *Env) (*Env, error) {
	local := make(map[string]LispExp)
	for params != nil {
		if args == nil {
			return nil, errors.New("mismatched params and args lengths")
		}
		param, ok := params.Car.(*SymbolAtom)
		if !ok {
			return nil, errors.New("parameters can only be symbols")
		}
		local[param.Data] = args.Car
		args = args.Cdr
		params = params.Cdr
	}
	return &Env{
		local: local,
		outer: outer,
	}, nil
}

func (e *Env) SetVar(name string, val LispExp) {
	e.local[name] = val
}

// find the innermost Env in which varName appears.
func (e *Env) FindEnvWith(varName string) *Env {
	for key := range e.local {
		if key == varName {
			return e
		}
	}
	if e.outer == nil {
		return nil
	}
	return e.outer.FindEnvWith(varName)
}

func (e *Env) FindVar(varName string) LispExp {
	envWithVar := e.FindEnvWith(varName)
	if envWithVar == nil {
		return nil
	}
	val, ok := envWithVar.local[varName]
	if !ok {
		return nil
	}
	return val
}

func standardEnv() *Env {
	params := NewList([]LispExp{
		&SymbolAtom{TRUE},
		&SymbolAtom{FALSE},
		&SymbolAtom{EQUAL},
		&SymbolAtom{ADD},
		&SymbolAtom{MULT},
		&SymbolAtom{DIV},
	}...)
	args := NewList([]LispExp{
		&SymbolAtom{TRUE},
		&SymbolAtom{FALSE},
		&BuiltIn{equalOp},
		&BuiltIn{addOp},
		&BuiltIn{multOp},
		&BuiltIn{divOp},
	}...)
	env, err := NewEnv(params, args, nil)
	if err != nil {
		panic("couldn't evaluate standard environment")
	}
	return env
}
