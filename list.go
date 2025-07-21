package main

// implementing the main List data structure as a linked list intitally because it
// is the most aligned with Lisp semantics. I do wonder how a slice  of arbitrary types
// might work though.

// consider renaming, and see whether back links are really necessary
type Linker interface {
	Next() Linker
	Prev() Linker
	SetNext(Linker)
	SetPrev(Linker)
}

type item struct {
	next Linker
	prev Linker
}

func (i *item) Next() Linker {
	return i.next
}

func (i *item) Prev() Linker {
	return i.prev
}

func (i *item) SetNext(next Linker) {
	i.next = next
}

func (i *item) SetPrev(prev Linker) {
	i.prev = prev
}

type IntItem struct {
	item
	Data int64
}

type FloatItem struct {
	item
	Data float64
}

type SymbolItem struct {
	item
	Data string
}

type ListItem struct {
	item
	Data *List
}

// Putting Lisp-like semantics over the Linker interface defined above
type List struct {
	head Linker
}

// head is a dummy head, an empty `item`
func NewList(items ...Linker) *List {
	head := item{}
	if len(items) == 0 {
		return &List{&head}
	}
	var current Linker = &head
	for _, item := range items {
		current.SetNext(item)
		item.SetPrev(current)
		current = item
	}
	return &List{&head}
}

func (l *List) Car() Linker {
	return l.head.Next()
}

func (l *List) Cdr() *List {
	next := l.Car().Next()
	return NewList(next)
}
