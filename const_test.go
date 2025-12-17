package main

import (
	"reflect"
	"testing"
)

func TestList(t *testing.T) {
	first := IntAtom{
		Data: 1,
	}
	second := IntAtom{
		Data: 2,
	}
	list := Cons(&first, Cons(&second, nil))
	firstNext := reflect.ValueOf(list.Cdr.Car).Elem().Interface()
	fnItem, ok := firstNext.(IntAtom)
	if !ok {
		t.Error("expected IntItem")
	}
	if fnItem.Data != second.Data {
		t.Errorf("got %d wanted %d", fnItem.Data, second.Data)
	}

	list = NewList(&first)
	outer := IntAtom{Data: 5}
	outerThird := SymbolAtom{Data: "x"}
	outerList := Cons(&outerThird, Cons(&outer, list))
	retrievedList := outerList.Cdr.Cdr
	if retrievedList.Car != &first {
		t.Errorf("got %v wanted %v", retrievedList.Car, &first)
	}
}
