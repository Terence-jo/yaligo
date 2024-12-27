package main

import (
	"regexp"
	"strconv"
)

// I need to define all types I'll use. First I need tokens, with
// a token type and a string value. Type will need to be based on
// regex right? Maybe not to start, Norvig just made an atom()
// function that tried casting the value as an int, then a float,
// then returned it as a Symbol. That's easy enough, the trick just
// becomes descending the nested program to add the atomized tokens
// to the list.

// Have a look at how the Go std lib implements tokens and parsing.
// That code (`go/token`, `go/parser`) handles Go parsing, but
// should generalise just fine. You can see how they handle different
// keywords and types of literals as `const` symbols. Start there.

type TokenClass int

type Pattern struct {
	tok    TokenClass
	regexp *regexp.Regexp
}

type Token struct {
	Class TokenClass
	Lit   string
}

const (
	OPEN TokenClass = iota
	CLOSE
	NUMBER
	SYMBOL
	DEFINE
	IF
)

// Create an array of strings, explicitly assigning items to each index
var tokens = [...]string{
	OPEN:   "(",
	CLOSE:  ")",
	NUMBER: "NUMBER",
	SYMBOL: "SYMBOL",
	DEFINE: "DEFINE",
	IF:     "IF",
}

func (t TokenClass) String() string {
	var s string
	if 0 <= t && t < TokenClass(len(tokens)) {
		s = tokens[t]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(t)) + ")"
	}
	return s
}

// Forgoing the full scanning functionality possessed by the
// std lib `go/scanner` package, just using regexp to determine
// the correct tokens.
var patterns = []Pattern{
	{OPEN, regexp.MustCompile(`^(\()`)},
	{CLOSE, regexp.MustCompile(`^(\))`)},
	{NUMBER, regexp.MustCompile(`^([0-9]+\.?[0-9]*)`)},
	{SYMBOL, regexp.MustCompile(`^('|[^\s();\.]+)`)},
}

// this might want to be renamed. it takes string tokens and replaces them with
// int/Token tokens, but it is distinct from the functionality in parse.go.
func LexTokens(programTokenised []string) []Token {
	var tokens []Token
	for _, lit := range programTokenised {
		for _, pattern := range patterns {
			matches := pattern.regexp.FindStringSubmatch(lit)
			if matches != nil {
				tokens = append(tokens, Token{pattern.tok, lit})
				break
			}
		}
	}
	return tokens
}
