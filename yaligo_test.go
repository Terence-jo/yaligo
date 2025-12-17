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
	intItem, ok := intAtom.(*IntAtom)
	if !ok {
		t.Errorf("failed to convert %v to IntItem", intAtom)
	}
	desiredInt := IntAtom{Data: 10}
	if !reflect.DeepEqual(*intItem, desiredInt) {
		t.Errorf("got %v, wanted %v", *intItem, desiredInt)
	}

	floatAtom, err := atom(floatTok)
	if err != nil {
		t.Errorf("invalid float atom %v", floatTok)
	}
	floatItem, ok := floatAtom.(*FloatAtom)
	if !ok {
		t.Errorf("failed to convert %v to FloatItem", floatAtom)
	}
	desiredFloat := FloatAtom{Data: 5.5}
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
		&IntAtom{Data: 5},
		&IntAtom{Data: 6},
	)
	referenceList := NewList(
		&SymbolAtom{Data: "define"},
		&SymbolAtom{Data: "x"},
		&IntAtom{Data: 10},
		&SymbolAtom{Data: "y"},
		innerList,
	)
	assertListEqual(t, parsedList, referenceList)
}

func TestEval(t *testing.T) {
	list := NewList(
		&SymbolAtom{Data: "+"},
		&IntAtom{Data: 1},
		&IntAtom{Data: 1},
	)
	got, err := Eval(list, globalEnv)
	if err != nil {
		t.Error(err)
	}
	want := 2.0
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
	eqProc := NewList(
		&SymbolAtom{Data: "eq?"},
		&IntAtom{Data: 1},
		&IntAtom{Data: 1},
	)
	ifList := NewList(
		&SymbolAtom{Data: "if"},
		eqProc,
		&FloatAtom{Data: 5.0},
		&FloatAtom{Data: 10.0},
	)
	got, err = Eval(ifList, globalEnv)
	if err != nil {
		t.Error(err)
	}
	want = 5.0
	if got != want {
		t.Errorf("got %f, wanted %f", got, want)
	}
	// test some failure cases
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
