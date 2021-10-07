// Package typedcontext implements type-preserving context functions.
package typedcontext

import (
	"context"
	"time"
)

// This can be put into a standalone library and used in different projects
// without modification.

// TypedContext is the interface required to make your context type
// type-preserved. T is typically the implementor type.
type TypedContext[T any] interface {
	context.Context
	InnerContext() context.Context
	WithInnerContext(context.Context) T
}

// WithTimeout returns sets the timeout on the underlying context in ctx and
// returns T.
func WithTimeout[T any](ctx TypedContext[T], timeout time.Duration) (T, context.CancelFunc) {
	inner, cancel := context.WithTimeout(ctx.InnerContext(), timeout)
	return ctx.WithInnerContext(inner), cancel
}
