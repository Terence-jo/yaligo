package main

import (
	"reflect"
	"testing"
)

func TestParseTokens(t *testing.T) {
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
	got := LexTokens(program)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %+v, wanted %+v", got, want)
	}
}
