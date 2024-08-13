package main

// I'm just going to use a blank interface for the
// car and rely on reflection initially. Later on I
// might pursue a more performant solution.
type List struct {
	car interface{}
	cdr *List
}
