package main

import (
	"reflect"
	"testing"
)

func TestTokenise(t *testing.T) {
	exp := "(define x 10)"
	want := []string{
		"(", "define", "x", "10", ")",
	}
	got := tokenise(exp)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestLexTokens(t *testing.T) {
	program := tokenise("(begin + (* 2.7 5) 6.4)")
	want := []Token{
		{OPEN, "("},
		{SYMBOL, "begin"},
		{SYMBOL, "+"},
		{OPEN, "("},
		{SYMBOL, "*"},
		{NUMBER, "2.7"},
		{NUMBER, "5"},
		{CLOSE, ")"},
		{NUMBER, "6.4"},
		{CLOSE, ")"},
	}
	got := lexTokens(program)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %+v, wanted %+v", got, want)
	}
}
