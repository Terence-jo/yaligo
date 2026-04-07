package internal

import (
	"testing"
)

type EnvExtension struct {
	Keys   *ConsCell
	Values *ConsCell
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
	if symbolItem.String() != desiredSymbol.String() {
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
	if intItem.String() != desiredNumber.String() {
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
	if floatItem.String() != desiredFloat.String() {
		t.Errorf("got %v, wanted %v", *floatItem, desiredFloat)
	}
}

func TestReadFromTokens(t *testing.T) {
	cases := []struct {
		exp      string
		expected *ConsCell
	}{
		{
			"(define x 10 y (- 5 6))",
			NewList(
				&SymbolAtom{Data: "define"},
				&SymbolAtom{Data: "x"},
				&NumberAtom{Data: 10},
				&SymbolAtom{Data: "y"},
				NewList(
					&SymbolAtom{Data: "-"},
					&NumberAtom{Data: 5},
					&NumberAtom{Data: 6},
				),
			),
		},
		{
			"(lambda (x) (x))",
			NewList(
				&SymbolAtom{"lambda"},
				NewList(&SymbolAtom{"x"}),
				NewList(&SymbolAtom{"x"}),
			),
		},
	}
	for _, tt := range cases {
		t.Run(tt.exp, func(t *testing.T) {
			toks := tokenise(tt.exp)
			lexed := lexTokens(toks)
			parsed, _, err := readFromTokens(lexed, 0)
			if err != nil {
				t.Error("failed to read tokens")
			}
			parsedList, ok := parsed.(*ConsCell)
			if !ok {
				t.Error("did not parse to a list")
			}
			assertListEqual(t, parsedList, tt.expected)
		})
	}
}

func TestEval(t *testing.T) {
	tests := []struct {
		name         string
		code         *ConsCell
		envAdditions EnvExtension
		expected     LispExp
	}{
		{
			"simpleAdd",
			NewList(
				&SymbolAtom{Data: ADD},
				&NumberAtom{Data: 1},
				&NumberAtom{Data: 1},
			),
			EnvExtension{nil, nil},
			&NumberAtom{2.0},
		},
		{
			"simpleEq",
			NewList(
				&SymbolAtom{Data: EQUAL},
				&NumberAtom{Data: 1},
				&NumberAtom{Data: 1},
			),
			EnvExtension{nil, nil},
			&SymbolAtom{TRUE},
		},
		{
			"falseEq",
			NewList(
				&SymbolAtom{Data: EQUAL},
				&NumberAtom{Data: 0},
				&NumberAtom{Data: 1},
			),
			EnvExtension{nil, nil},
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
			EnvExtension{nil, nil},
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
			EnvExtension{nil, nil},
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
			EnvExtension{nil, nil},
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
			EnvExtension{nil, nil},
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
		{
			"procCall",
			NewList(
				&SymbolAtom{"add-2"},
				&NumberAtom{2},
			),
			EnvExtension{
				&ConsCell{
					&SymbolAtom{"add-2"},
					nil,
				},
				&ConsCell{
					&Procedure{
						params: NewList(&SymbolAtom{"x"}),
						body: NewList(
							&SymbolAtom{ADD},
							&SymbolAtom{"x"},
							&NumberAtom{2},
						),
						env: globalEnv,
					},
					nil,
				},
			},
			&NumberAtom{4},
		},
		{
			"quote",
			NewList(
				&SymbolAtom{"quote"},
				NewList(&SymbolAtom{"quoted"}, &SymbolAtom{"list"}),
			),
			EnvExtension{nil, nil},
			NewList(&SymbolAtom{"quoted"}, &SymbolAtom{"list"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var env *Env
			var err error
			if tt.envAdditions.Keys == nil {
				env = globalEnv
			} else {
				env, err = NewEnv(tt.envAdditions.Keys, tt.envAdditions.Values, globalEnv)
				if err != nil {
					t.Error(err)
				}
			}
			got, err := Eval(tt.code, env)
			if err != nil {
				t.Error(err)
			}
			if got.String() != tt.expected.String() {
				t.Errorf("got %v, wanted %v", got, tt.expected)
			}
		})
	}
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
		if testNext != nil && refNext == nil {
			t.Fatalf("mismatch between testList cdr and reference cdr with %v and nil", testNext)
		}
		if refNext != nil && testNext == nil {
			t.Fatalf("mismatch between testList cdr and reference cdr with %v and nil", refNext)
		}
		if testNext == nil && refNext == nil {
			return
		}
		testList = testNext
		referenceList = refNext
	}
}

func assertListIter(t testing.TB, testList *ConsCell, referenceList *ConsCell) {
	t.Helper()
	got := testList.Car.String()
	want := referenceList.Car.String()
	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}
}
