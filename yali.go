package main

import (
	"errors"
	"strconv"
	"strings"
)

var globalEnv = standardEnv()

func Eval(exp LispExp, env *Env) (LispExp, error) {
	switch exp := exp.(type) {
	case *NumberAtom:
		return exp, nil
	case *SymbolAtom:
		symbol := exp.Data
		return env.FindVar(symbol), nil
	case *ConsCell:
		car := exp.Car
		exp = exp.Cdr
		if car == nil {
			return &ConsCell{nil, nil}, nil
		}
		// unquoted list, this means it is a syntactic form or procedure call:
		symbol, ok := car.(*SymbolAtom)
		if !ok {
			return nil, errors.New("expected symbol at head of unquoted list")
		}
		switch symbol.Data {
		case "quote":
			return exp.Car, nil
		case "if":
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
			resSym, ok := testRes.(*SymbolAtom)
			if !ok {
				return nil, errors.New("expected test of 'if' to evaluate to a boolean")
			}
			if resSym.Data == TRUE {
				return Eval(conseq, env)
			}
			if resSym.Data != FALSE {
				return nil, errors.New("expected test of 'if' to evaluate to a boolean")
			}
			return Eval(alt, env)
		case "cond":
			for exp != nil {
				clause, ok := exp.Car.(*ConsCell)
				if !ok {
					return nil, errors.New("expect cond clauses to be cons cells")
				}
				if clause.Car == nil || clause.Cdr == nil {
					return nil, errors.New("cond clauses must have two elements")
				}
				testRes, err := Eval(clause.Car, env)
				if err != nil {
					return nil, err
				}
				resSym, ok := testRes.(*SymbolAtom)
				if !ok {
					return nil, errors.New("expected test to evaluate to boolean literal")
				}
				if resSym.Data == TRUE {
					return Eval(clause.Cdr.Car, env)
				}
				if resSym.Data != FALSE {
					return nil, errors.New("expected test to evaluate to boolean literal")
				}
				exp = exp.Cdr
			}
			return &ConsCell{nil, nil}, nil
		case "define":
			varName, ok := exp.Car.(*SymbolAtom)
			if !ok {
				return nil, errors.New("expected a symbol as first arg to 'define'")
			}
			expArg, err := Eval(exp.Cdr.Car, env)
			if err != nil {
				return nil, err
			}
			env.SetVar(varName.Data, expArg)
			return expArg, nil
		case "set!":
			varName, ok := exp.Car.(*SymbolAtom)
			if !ok {
				return nil, errors.New("expected a symbol as first arg to 'define'")
			}
			expArg := exp.Cdr.Car
			envWithVar := env.FindEnvWith(varName.Data)
			if envWithVar == nil {
				return nil, errors.New("expected target of 'set' to exist")
			}
			envWithVar.SetVar(varName.Data, expArg)
			return expArg, nil
		case "lambda":
			// gather params, body, env and create a Procedure
			params, ok := exp.Car.(*ConsCell)
			if !ok {
				return nil, errors.New("expected list of symbols for params")
			}
			return &Procedure{params, exp.Cdr.Car, env}, nil
		default:
			procExp, err := Eval(symbol, env)
			if err != nil {
				return nil, err
			}
			proc, ok := procExp.(Callable)
			if !ok {
				return nil, errors.New("expected procedure name at head of list")
			}

			var evalArgs func(args *ConsCell) (*ConsCell, error)
			evalArgs = func(args *ConsCell) (*ConsCell, error) {
				if args == nil {
					return nil, nil
				}
				argVal, err := Eval(args.Car, env)
				if err != nil {
					return nil, err
				}
				tail, err := evalArgs(args.Cdr)
				if err != nil {
					return nil, err
				}
				return Cons(argVal, tail), nil
			}

			argVals, err := evalArgs(exp)
			if err != nil {
				return nil, err
			}
			ret, err := proc.Call(argVals)
			if err != nil {
				return nil, err
			}
			return ret, nil
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
		// handle empty list '()'
		if tokens[pos].Class == CLOSE {
			return &ConsCell{nil, nil}, pos + 1, nil
		}

		// read the first element to create the head of the list.
		car, newPos, err := readFromTokens(tokens, pos)
		if err != nil {
			return nil, 0, err
		}
		head := Cons(car, nil)
		tail := head
		pos = newPos

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
		pos++
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
		intval, err := strconv.ParseInt(token.Lit, 10, 64)
		if err != nil {
			floatval, err := strconv.ParseFloat(token.Lit, 64)
			if err != nil {
				return nil, err
			}
			return &NumberAtom{Data: floatval}, nil
		}
		return &NumberAtom{Data: float64(intval)}, nil
	}
	return &SymbolAtom{Data: token.Lit}, nil
}
