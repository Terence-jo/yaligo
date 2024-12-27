package main

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
	*list
}

// Putting Lisp-like semantics over the Linker interface defined above
type list struct {
	head Linker
}

// head is a dummy head, an empty `item`
func List(firstItem Linker) *list {
	head := item{}
	head.SetNext(firstItem)
	firstItem.SetPrev(&head)
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
