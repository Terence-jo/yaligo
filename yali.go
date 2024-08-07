package main

import "strings"

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

// This is going to require quite a different approach to lis.py
// due to Go's type rigidity. I'll need to be dirty at the low
// level and make use of {}interface, and I will need to define
// value and cons interfaces to allow the seamless processing
// of any part of the expression.
func readFromTokens(tokens []string) []
