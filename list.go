package main

// implementing the main List data structure as a linked list intitally because it
// is the most aligned with Lisp semantics. I do wonder how a slice  of arbitrary types
// might work though.

// consider renaming, and see whether back links are really necessary
type LispExp interface {
	Next() LispExp
	SetNext(LispExp)
}

type item struct {
	next LispExp
}

func (i *item) Next() LispExp {
	return i.next
}

func (i *item) SetNext(next LispExp) {
	i.next = next
}

type Atom interface {
	Value() any
}

type IntAtom struct {
	item
	data int64
}

func (i *IntAtom) Value() any {
	return i.data
}

type FloatAtom struct {
	item
	data float64
}

func (f *FloatAtom) Value() any {
	return f.data
}

type SymbolAtom struct {
	item
	data string
}

func (s *SymbolAtom) Value() any {
	return s.data
}

type ListExp struct {
	item
}

func NewList(items ...LispExp) *ListExp {
	exp := ListExp{item{}}
	if len(items) == 0 {
		return &exp
	}
	var current LispExp = &exp
	for _, item := range items {
		current.SetNext(item)
		current = item
	}
	return &exp
}

func (l *ListExp) Car() LispExp {
	return l.Next()
}

func (l *ListExp) Cdr() *ListExp {
	head := l.Car()
	if head.Next() == nil {
		return nil
	}
	return NewList(head.Next())
}
