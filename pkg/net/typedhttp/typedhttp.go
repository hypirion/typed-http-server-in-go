// Package typedhttp provides typed HTTP types via generics.
package typedhttp

import (
	"context"
	"net/http"
)

// My guess is that we'll have these standardised types that people implement
// decoders and encoders for, but with a specific error type that you can attach
// status codes (and possibly more) on.

// A Handler handles a request-response request where In is the input and Out/an
// error is the output.
type Handler[Ctx context.Context, In, Out any] func(Ctx, In, *http.Request) (Out, error)


