package main

// This file implements Phase 1 of the cons cell refactor plan.
// It defines the new data model in a separate file to avoid editing existing code.

// LispExp is the interface for all Lisp expressions in the new model.
// For the purpose of this isolated implementation, we'll use an empty interface
// to represent any Lisp value. In a full refactor, this might be more specific.
type LispExp interface {
	isLispExp()
}

// ConsCell is the new primary list structure, representing a standard Lisp cons cell.
type ConsCell struct {
	Car LispExp
	Cdr LispExp
}

func (c *ConsCell) isLispExp() {}

// IntAtom represents an integer atom. It no longer embeds the old 'item' struct.
type IntAtom struct {
	Data int64
}

func (i *IntAtom) isLispExp() {}

// FloatAtom represents a float atom. It no longer embeds the old 'item' struct.
type FloatAtom struct {
	Data float64
}

func (f *FloatAtom) isLispExp() {}

// SymbolAtom represents a symbol atom. It no longer embeds the old 'item' struct.
type SymbolAtom struct {
	Data string
}

func (s *SymbolAtom) isLispExp() {}

// Cons creates a new ConsCell (a new list node).
func Cons(car LispExp, cdr *ConsCell) *ConsCell {
	return &ConsCell{Car: car, Cdr: cdr}
}

// NewList creates a new proper list (a nil-terminated chain of ConsCells)
// from a slice of LispExp items.
func NewList(items ...LispExp) LispExp {
	if len(items) == 0 {
		return nil // Represents the empty list '()'.
	}
	// Build the list backwards from the last element to the first.
	var list *ConsCell = nil
	for _, item := range items {
		list = Cons(item, list)
	}
	return list
}
