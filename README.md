# stacker

[![GoDoc](https://godoc.org/github.com/olefasting/stacker?status.svg)](https://godoc.org/github.com/olefasting/stacker) 
[![GoCover](http://gocover.io/_badge/github.com/olefasting/stacker)](http://gocover.io/github.com/olefasing/stacker)
[![Coveralls](https://coveralls.io/repos/github/olefasting/stacker/badge.svg?branch=master)](https://coveralls.io/github/olefasting/stacker?branch=master)
[![Build Status](https://travis-ci.org/olefasting/stacker.svg?branch=master)](https://travis-ci.org/olefasting/stacker)

Package **stacker** provides middleware stacking for [httpctx](http://github.com/olefasting/httpctx) handlers.
This is more or less a rewrite of [justinas/alice](https://github.com/justinas/alice), so it works pretty much the same, except for the removal of the `Extend` and the `ThenFunc` methods.
Since `Stack` is type `[]Middleware`, you can just pass it as an argument to `Append`, and it will be unpacked automatically.

## Documentation

Create new stacks by calling `New`, add middleware by calling `Append`. The stacks are **immutable**, so you can create stacks for various purposes and mix and match at the end points.

This depends on [httpctx](http://github.com/olefasting/httpctx)

#### Example:

```go
package main

import (
	"net/http"

	"golang.org/x/net/context"

	"github.com/olefasting/httpctx"
	"github.com/olefasting/stacker"
)

func main() {
	// App handler
	ah := &appHandler{}

	// Create stack with optional arguments
	stack1 := stacker.New(mw1, mw2)

	// Append some more middleware.
	// Remember that the stack is immutable, so you
	// can either save the new stack over the old,
	// or you can create a new stack
	stack2 := stack1.Append(mw1, mw2)

	// You can also use stacks as arguments when
	// calling Append, just unpack it like here
	stack2 = stack2.Append(stack1...)

	// Then will create a handler that starts with
	// the first middleware handler in the stack, and
	// ends with the ah handler
	h := stack2.Then(ah)

	// Adapt for net/http ServeMux
	ctx := context.Background()
	app := httpctx.AdaptHandler(ctx, h)

	// Listen and serve
	http.ListenAndServe(":8080", app)
}

// Handler implementation
type appHandler struct{}

// Handler method
func (h *appHandler) ServeHTTPCtx(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	// Get something from context
	body := ctx.Value("key").(string)

	// App does it stuff

	// Write
	rw.Write([]byte(body))
}

// Middleware1
func mw1(next httpctx.Handler) httpctx.Handler {
	fn := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		// Add value to context
		newctx := context.WithValue(ctx, "key", "some thing")

		// Pass to next
		next.ServeHTTPCtx(newctx, rw, req)
	}

	// Return
	return httpctx.HandlerFunc(fn)
}

// Middleware2
func mw2(next httpctx.Handler) httpctx.Handler {
	fn := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		// Add value to context
		newctx := context.WithValue(ctx, "key", "some other thing")

		// Pass to next
		next.ServeHTTPCtx(newctx, rw, req)
	}

	// Return
	return httpctx.HandlerFunc(fn)
}

```

Licence MIT
