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
	return
}

func (i *item) SetPrev(prev Linker) {
	i.prev = prev
	return
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
	Data *list
}

// Putting Lisp-like semantics over the Linker interface defined above
type list struct {
	head Linker
}

// head is a dummy head, an empty `item`
func List(items ...Linker) *list {
	head := item{}
	if len(items) == 0 {
		return &list{&head}
	}
	var cur Linker = &head
	for _, item := range items {
		cur.SetNext(item)
		item.SetPrev(cur)
		cur = item
	}
	return &list{&head}
}

func (l *list) Car() Linker {
	return l.head.Next()
}

func (l *list) Cdr() *list {
	next := l.Car().Next()
	return List(next)
}

// func (l *List) Eval()
