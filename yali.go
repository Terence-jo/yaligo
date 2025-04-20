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

func readFromTokens(tokens []Token, pos int) (Linker, int, error) {
	if len(tokens) == 0 {
		return nil, 0, errors.New("unexpected EOF")
	}
	token := tokens[pos]
	if token.Class == OPEN {
		// start a list and recur to fill it
		pos++
		list := List()
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
	} else if token.Class == CLOSE {
		return nil, 0, errors.New("unexpected )")
	}
	atom, err := atom(token)
	if err != nil {
		return nil, 0, err
	}
	pos++
	return atom, pos, nil
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
