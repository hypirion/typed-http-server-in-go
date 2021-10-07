package plusexample

import (
	"context"
	"fmt"
)

// This is the "tying the knot" example from the blogpost
// https://hypirion.com/musings/type-safe-http-servers-in-go-via-generics
//
// It has nothing to do with the typed HTTP server itself

type Plusser[T any] interface {
	Plus(T) T
}

type Int int

func (a Int) Plus(b Int) Int {
	return a + b
}

var _ Plusser[Int] = Int(0)

//           this one ↓
func Concat[T Plusser[T]](xs []T) T {
  val := xs[0] // panic if slice is empty
  for _, x := range xs[1:] {
    val = val.Plus(x)
  }
  return val
}

func DoSomething1(ctx context.Context) {}

func DoSomething2[T context.Context](ctx T) {}

func init() {
	// Look, they have the same type signature:
	x := DoSomething1
	x = DoSomething2[context.Context]
	// But if I tried to make the type anything else, it wouldn't work.

	// (To avoid compiler errors)
	if false {
		fmt.Println(x)
	}
}
