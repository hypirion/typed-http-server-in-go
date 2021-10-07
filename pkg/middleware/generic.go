// Package middleware implements middleware functions. It contains both
// app-specific middleware functions, as well as more generic ones.
package middleware

import (
	"net/http"
	"time"

	tc "github.com/hypirion/typed-http-server-in-go/pkg/typedcontext"
	th "github.com/hypirion/typed-http-server-in-go/pkg/net/typedhttp"
)

// SetHandlerTimeout is so generic that you can probably put it into a generic
// middleware handler library (third party code).

// SetHandlerTimeout returns a new handler with its timeout set to timeout.
func SetHandlerTimeout[Ctx tc.TypedContext[Ctx], In, Out any](timeout time.Duration, h th.Handler[Ctx, In, Out]) th.Handler[Ctx, In, Out]{
	return func(ctx Ctx, in In, r *http.Request) (Out, error) {
		newCtx, cancel := tc.WithTimeout[Ctx](ctx, timeout)
		defer cancel()
		return h(newCtx, in, r)
	}
}
