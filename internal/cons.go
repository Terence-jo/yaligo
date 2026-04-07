package internal

import "fmt"

type LispExp interface {
	fmt.Stringer
	isLispExp()
}

type ConsCell struct {
	Car LispExp
	Cdr *ConsCell
}

func (c *ConsCell) isLispExp() {}

func (c *ConsCell) String() string {
	return fmt.Sprintf("(%v %v)", c.Car, c.Cdr)
}

type NumberAtom struct {
	Data float64
}

func (n *NumberAtom) isLispExp() {}

func (n *NumberAtom) String() string {
	return fmt.Sprintf("%f", n.Data)
}

type SymbolAtom struct {
	Data string
}

func (s *SymbolAtom) isLispExp() {}

func (s *SymbolAtom) String() string {
	return fmt.Sprintf("%s", s.Data)
}

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
