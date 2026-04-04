package main

import (
	"reflect"
	"testing"
)

func TestFindVar(t *testing.T) {
	env := standardEnv()
	// find in localEnv
	got := env.FindVar(TRUE)
	want := &SymbolAtom{TRUE}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
	// find in outer env
	// innerEnv, err := NewEnv(nil, nil, standardEnv())
	// if err != nil {
	// 	t.Error(err)
	// }
	innerEnv := &Env{
		map[string]LispExp{},
		standardEnv(),
	}
	got = innerEnv.FindVar(TRUE)
	want = &SymbolAtom{TRUE}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

}
