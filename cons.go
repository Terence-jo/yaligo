package main

type LispExp interface {
	isLispExp()
}

type ConsCell struct {
	Car LispExp
	Cdr *ConsCell
}

func (c *ConsCell) isLispExp() {}

type IntAtom struct {
	Data int64
}

func (i *IntAtom) isLispExp() {}

type FloatAtom struct {
	Data float64
}

func (f *FloatAtom) isLispExp() {}

type SymbolAtom struct {
	Data string
}

func (s *SymbolAtom) isLispExp() {}

func Cons(car LispExp, cdr *ConsCell) *ConsCell {
	return &ConsCell{Car: car, Cdr: cdr}
}

func NewList(items ...LispExp) *ConsCell {
	if len(items) == 0 {
		return nil // Represents the empty list '()'.
	}
	// Build the list backwards from the last element to the first.
	var list *ConsCell = nil
	for i := len(items) - 1; i >= 0; i-- {
		list = Cons(items[i], list)
	}
	return list
}
