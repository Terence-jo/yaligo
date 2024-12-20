package main

import (
	"reflect"
	"testing"
)

func TestList(t *testing.T) {
	first := IntItem{
		// Item: Item{},
		Data: 1,
	}
	second := IntItem{
		// Item: Item{},
		Data: 2,
	}
	first.SetNext(&second)
	second.SetPrev(&first)
	firstNext := reflect.ValueOf(first.Next()).Elem().Interface()
	fnItem, ok := firstNext.(*IntItem)
	if !ok {
		t.Error("expected IntItem")
	}
	if fnItem.Data != second.Data {
		t.Errorf("got %d wanted %d", fnItem.Data, second.Data)
	}

	list := ListItem{
		List: List{&first},
	}
	outer := IntItem{Data: 5}
	outerThird := SymbolItem{Data: "x"}
	outer.SetNext(&list)
	list.SetPrev(&outer)
	list.SetNext(&outerThird)
	outerThird.SetPrev(&list)
	retrievedList, ok := reflect.ValueOf(outer.Next()).Elem().Interface().(*ListItem)
	if !ok {
		t.Error("expected a ListItem")
	}
	if retrievedList.Head != &first {
		t.Errorf("got %v wanted %v", retrievedList.Head, &first)
	}
}
