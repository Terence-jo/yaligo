package main

// I need to define all types I'll use. First I need tokens, with
// a token type and a string value. Type will need to be based on
// regex right? Maybe not to start, Norvig just made an atom()
// function that tried casting the value as an int, then a float,
// then returned it as a Symbol. That's easy enough, the trick just
// becomes descending the nested program to add the atomized tokens
// to the list.

// Have a look at how the Go std lib implements tokens and parsing.
// That code (`go/token`, `go/parser`) handles Go parsing, but
// should generalise just fine. You can see how they handle different
// keywords and types of literals as `const` symbols. Start there.

type Atom interface {
	Read(string) interface{}
}
type Symbol string

type List struct {
	car ConsValue
	cdr ConsValue
}
