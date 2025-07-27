package main

import (
	"reflect"
	"testing"
)

func TestList(t *testing.T) {
	first := IntAtom{
		data: 1,
	}
	second := IntAtom{
		data: 2,
	}
	first.SetNext(&second)
	second.SetPrev(&first)
	firstNext := reflect.ValueOf(first.Next()).Elem().Interface()
	fnItem, ok := firstNext.(IntAtom)
	if !ok {
		t.Error("expected IntItem")
	}
	if fnItem.data != second.data {
		t.Errorf("got %d wanted %d", fnItem.data, second.data)
	}

	list := NewList(&first)
	outer := IntAtom{data: 5}
	outerThird := SymbolAtom{data: "x"}
	outer.SetNext(list)
	list.SetPrev(&outer)
	list.SetNext(&outerThird)
	outerThird.SetPrev(list)
	retrievedList, ok := reflect.ValueOf(outer.Next()).Elem().Interface().(ListExp)
	if !ok {
		t.Error("expected a ListItem")
	}
	if retrievedList.Car() != &first {
		t.Errorf("got %v wanted %v", retrievedList.Car(), &first)
	}
}
