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
	symbolItem, ok := symbolAtom.(*SymbolItem)
	if !ok {
		t.Errorf("failed to convert %v to SymbolItem", symbolAtom)
	}
	desiredSymbol := SymbolItem{Data: "define"}
	if !reflect.DeepEqual(*symbolItem, desiredSymbol) {
		t.Errorf("got %v, wanted %v", *symbolItem, desiredSymbol)
	}

	intAtom, err := atom(intTok)
	if err != nil {
		t.Errorf("invalid int atom %v", intTok)
	}
	intItem, ok := intAtom.(*IntItem)
	if !ok {
		t.Errorf("failed to convert %v to IntItem", intAtom)
	}
	desiredInt := IntItem{Data: 10}
	if !reflect.DeepEqual(*intItem, desiredInt) {
		t.Errorf("got %v, wanted %v", *intItem, desiredInt)
	}

	floatAtom, err := atom(floatTok)
	if err != nil {
		t.Errorf("invalid float atom %v", floatTok)
	}
	floatItem, ok := floatAtom.(*FloatItem)
	if !ok {
		t.Errorf("failed to convert %v to FloatItem", floatAtom)
	}
	desiredFloat := FloatItem{Data: 5.5}
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
	parsedList, ok := parsed.(*ListItem)
	if !ok {
		t.Error("did not parse to a list")
	}
	innerList := List(
		&SymbolItem{Data: "-"},
		&IntItem{Data: 5},
		&IntItem{Data: 6},
	)
	referenceList := &ListItem{
		Data: List(
			&SymbolItem{Data: "define"},
			&SymbolItem{Data: "x"},
			&IntItem{Data: 10},
			&SymbolItem{Data: "y"},
			&ListItem{Data: innerList},
		),
	}
	assertListEqual(t, parsedList, referenceList)
	// if !reflect.DeepEqual(parsed, list) {
	// 	t.Errorf("got %v, wanted %v", parsed, list)
	// }
}

func assertListEqual(t testing.TB, testList *ListItem, referenceList *ListItem) {
	for true {
		switch testList.Data.Car().(type) {
		case *ListItem:
			innerTestList := testList.Data.Car().(*ListItem)
			innerReferenceList := referenceList.Data.Car().(*ListItem)
			assertListEqual(t, innerTestList, innerReferenceList)
		default:
			assertListIter(t, testList, referenceList)
		}
		testNext := testList.Data.Car().Next()
		refNext := referenceList.Data.Car().Next()
		if testNext != nil {
			if refNext == nil {
				t.Errorf("mismatch between testList cdr and reference cdr with %v and nil", testNext)
			}
			testList.Data = testList.Data.Cdr()
			referenceList.Data = referenceList.Data.Cdr()
		} else {
			return
		}
	}
}

func assertListIter(t testing.TB, testList *ListItem, referenceList *ListItem) {
	got := reflect.Indirect(reflect.ValueOf(testList.Data.Car())).Field(1)
	want := reflect.Indirect(reflect.ValueOf(referenceList.Data.Car())).Field(1)
	if !got.Equal(want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
