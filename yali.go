package main

import (
	"errors"
	"strconv"
	"strings"
)

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
func Eval(l *List) (any, error) {
	// We have a list
	return 1, nil
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
