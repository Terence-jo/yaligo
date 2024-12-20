package main

type Linker interface {
	Next() *Linker
	Prev() *Linker
	SetNext(Linker)
	SetPrev(Linker)
}

type item struct {
	next *Linker
	prev *Linker
}

func (i *item) Next() *Linker {
	return i.next
}

func (i *item) Prev() *Linker {
	return i.prev
}

func (i *item) SetNext(next Linker) {
	i.next = &next
	return
}

func (i *item) SetPrev(prev Linker) {
	i.prev = &prev
	return
}

type IntItem struct {
	item
	Data int
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
	List
}

// This was a little misguided, revisit
type List struct {
	Head Linker
}

// func (l *List) Eval()
