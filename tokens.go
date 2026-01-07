package main

import (
	"regexp"
	"strconv"
)

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
