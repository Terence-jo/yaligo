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

func TestAtom(t *testing.T) {
	symbolTok := Token{Class: SYMBOL, Lit: "define"}
	intTok := Token{Class: NUMBER, Lit: "10"}
	floatTok := Token{Class: NUMBER, Lit: "5.5"}

	symbolAtom, err := atom(symbolTok)
	if err != nil {
		t.Errorf("invalid symbol atom %v", symbolTok)
	}
	symbolItem, ok := symbolAtom.(*SymbolAtom)
	if !ok {
		t.Errorf("failed to convert %v to SymbolItem", symbolAtom)
	}
	desiredSymbol := SymbolAtom{Data: "define"}
	if !reflect.DeepEqual(*symbolItem, desiredSymbol) {
		t.Errorf("got %v, wanted %v", *symbolItem, desiredSymbol)
	}

	intAtom, err := atom(intTok)
	if err != nil {
		t.Errorf("invalid int atom %v", intTok)
	}
	intItem, ok := intAtom.(*NumberAtom)
	if !ok {
		t.Errorf("failed to convert %v to NumberItem", intAtom)
	}
	desiredNumber := NumberAtom{Data: 10}
	if !reflect.DeepEqual(*intItem, desiredNumber) {
		t.Errorf("got %v, wanted %v", *intItem, desiredNumber)
	}

	floatAtom, err := atom(floatTok)
	if err != nil {
		t.Errorf("invalid float atom %v", floatTok)
	}
	floatItem, ok := floatAtom.(*NumberAtom)
	if !ok {
		t.Errorf("failed to convert %v to FloatItem", floatAtom)
	}
	desiredFloat := NumberAtom{Data: 5.5}
	if !reflect.DeepEqual(*floatItem, desiredFloat) {
		t.Errorf("got %v, wanted %v", *floatItem, desiredFloat)
	}
}

func TestReadFromTokens(t *testing.T) {
	// These two lines are already tested above, trust them to work
	exp := "(define x 10 y (- 5 6))"
	// exp := "(define x 10 y 20)"
	toks := tokenise(exp)
	lexed := LexTokens(toks)
	parsed, _, err := readFromTokens(lexed, 0)
	if err != nil {
		t.Error("failed to read tokens")
	}
	parsedList, ok := parsed.(*ConsCell)
	if !ok {
		t.Error("did not parse to a list")
	}
	innerList := NewList(
		&SymbolAtom{Data: "-"},
		&NumberAtom{Data: 5},
		&NumberAtom{Data: 6},
	)
	referenceList := NewList(
		&SymbolAtom{Data: "define"},
		&SymbolAtom{Data: "x"},
		&NumberAtom{Data: 10},
		&SymbolAtom{Data: "y"},
		innerList,
	)
	assertListEqual(t, parsedList, referenceList)
}

func TestEval(t *testing.T) {
	tests := []struct {
		Name   string
		Code   *ConsCell
		Result LispExp
	}{
		{
			"simpleAdd",
			NewList(
				&SymbolAtom{Data: ADD},
				&NumberAtom{Data: 1},
				&NumberAtom{Data: 1},
			),
			&NumberAtom{2.0},
		},
		{
			"simpleEq",
			NewList(
				&SymbolAtom{Data: EQUAL},
				&NumberAtom{Data: 1},
				&NumberAtom{Data: 1},
			),
			&SymbolAtom{TRUE},
		},
		{
			"falseEq",
			NewList(
				&SymbolAtom{Data: EQUAL},
				&NumberAtom{Data: 0},
				&NumberAtom{Data: 1},
			),
			&SymbolAtom{FALSE},
		},
		{
			"trueIf",
			NewList(
				&SymbolAtom{Data: "if"},
				NewList(
					&SymbolAtom{Data: EQUAL},
					&NumberAtom{Data: 1},
					&NumberAtom{Data: 1},
				),
				&NumberAtom{Data: 5.0},
				&NumberAtom{Data: 10.0},
			),
			&NumberAtom{5.0},
		},
		{
			"falseIf",
			NewList(
				&SymbolAtom{Data: "if"},
				NewList(
					&SymbolAtom{Data: EQUAL},
					&NumberAtom{Data: 0},
					&NumberAtom{Data: 1},
				),
				&NumberAtom{Data: 5.0},
				&NumberAtom{Data: 10.0},
			),
			&NumberAtom{10.0},
		},
		{
			"cond",
			NewList(
				&SymbolAtom{Data: "cond"},
				NewList(
					NewList(
						&SymbolAtom{Data: EQUAL},
						&NumberAtom{Data: 0},
						&NumberAtom{Data: 1},
					),
					&NumberAtom{Data: 5.0},
				),
				NewList(
					NewList(
						&SymbolAtom{Data: EQUAL},
						&NumberAtom{Data: 1},
						&NumberAtom{Data: 1},
					),
					&NumberAtom{Data: 10.0},
				),
			),
			&NumberAtom{10.0},
		},
		{
			"lambda",
			NewList(
				&SymbolAtom{Data: "lambda"},
				NewList(
					&SymbolAtom{Data: "x"},
				),
				NewList(
					&SymbolAtom{"+"},
					&SymbolAtom{"x"},
					&NumberAtom{2},
				),
			),
			&Procedure{
				params: NewList(&SymbolAtom{"x"}),
				body: NewList(
					&SymbolAtom{"+"},
					&SymbolAtom{"x"},
					&NumberAtom{2},
				),
				env: globalEnv,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got, err := Eval(tt.Code, globalEnv)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, tt.Result) {
				t.Errorf("got %v, wanted %v", got, tt.Result)
			}
		})
	}
	// test some failure cases
	// test lambda
}

func assertListEqual(t testing.TB, testList *ConsCell, referenceList *ConsCell) {
	t.Helper()
	for {
		switch car := testList.Car.(type) {
		case *ConsCell:
			innerReferenceList := referenceList.Car.(*ConsCell)
			assertListEqual(t, car, innerReferenceList)
		default:
			assertListIter(t, testList, referenceList)
		}
		testNext := testList.Cdr
		refNext := referenceList.Cdr
		if testNext != nil {
			if refNext == nil {
				t.Errorf("mismatch between testList cdr and reference cdr with %v and nil", testNext)
			}
			testList = testNext
			referenceList = refNext
		} else {
			return
		}
	}
}

func assertListIter(t testing.TB, testList *ConsCell, referenceList *ConsCell) {
	t.Helper()
	got := reflect.Indirect(reflect.ValueOf(testList.Car)).Field(0)
	want := reflect.Indirect(reflect.ValueOf(referenceList.Car)).Field(0)
	if !got.Equal(want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
