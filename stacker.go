// Package stacker provides middleware stacking for http handlers
package stacker

import "net/http"

// Middleware is a wrapper for handlers that can be stacked
type Middleware func(http.Handler) http.Handler

// Stack is an immutable slice of Middleware handlers and the
// methods to manipulate  and create a handler from it.
type Stack []Middleware

// New will create a new handler stack
func New(mws ...Middleware) Stack {
	return Stack(mws)
}

// Append Middleware to a copy of this stack and return copy.
//
// By passing Stack... as the last argument to this function,
// you can append another stack to this.
func (s Stack) Append(mws ...Middleware) Stack {
	return append(s, mws...)
}

// Then will create a new handler from the Middleware in the
// stack, setting l as the last handler to be called.
//
// If nil is passed as l, nil is returned.
func (s Stack) Then(l http.Handler) http.Handler {
	// Iterate stack
	for i := len(s) - 1; i >= 0; i-- {
		l = s[i](l)
	}

	// Return finished handler
	return l
}
