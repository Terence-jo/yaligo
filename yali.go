package main

import (
	"errors"
	"strconv"
	"strings"
)

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

func readFromTokens(tokens []Token) (Linker, error) {
	if len(tokens) == 0 {
		return nil, errors.New("unexpected EOF")
	}
	token := tokens[0]
	tokens = tokens[1:]
	// need to think harder about how this should work and testing it.
	// to make it more general I might need a dummy head in the List
	// and have Car() return its next element, so I would access it
	// directly here.
	if token.Class == OPEN {
		// start a list and recur to fill it
		list := List(nil)
		cur := list.head
		for tokens[0].Class != CLOSE {
			next, err := readFromTokens(tokens)
			if err != nil {
				return nil, err
			}
			cur.SetNext(next)
			cur = next
		}
		expr := ListItem{list: list}
		return &expr, nil
	} else if token.Class == CLOSE {
		return nil, errors.New("unexpected )")
	}
	return atom(token)
}

func atom(token Token) (Linker, error) {
	if token.Class == NUMBER {
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
