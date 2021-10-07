package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	th "github.com/hypirion/typed-http-server-in-go/pkg/net/typedhttp"
	tc "github.com/hypirion/typed-http-server-in-go/pkg/typedcontext"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// I've joined packages together for brevity. This piece of the middleware
// package is probably something you define yourself in the source code of your
// program, whereas the other is something that can be defined in a generic
// middleware repository.

// UserID represents the ID of the user. You'd probably have UserID somewhere
// else than in the middleware package though.
type UserID int64

// UserIDFromString translates s into an UserID if it is valid.
func UserIDFromString(s string) (UserID, error) {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID: %w", err)
	}
	return UserID(val), err
}

// MyAppContext contains the UserID and a request-specific logger.
type MyAppContext struct {
	context.Context
	UserID UserID
	Logger logrus.FieldLogger
}

var _ context.Context = MyAppContext{}

// InnerContext implements the typedcontext.TypedContext interface.
func (mac MyAppContext) InnerContext() context.Context {
	return mac.Context
}

// WithInnerContext implements the typedcontext.TypedContext interface.
func (mac MyAppContext) WithInnerContext(ctx context.Context) MyAppContext {
	mac.Context = ctx
	return mac
}

var _ tc.TypedContext[MyAppContext] = MyAppContext{}

// AttachUserAndLogger fetches and attaches the user ID and attaches a
// request-specific logger.
func AttachUserAndLogger[Ctx context.Context, In, Out any](h th.Handler[MyAppContext, In, Out]) th.Handler[Ctx, In, Out] {
	return func(ctx Ctx, in In, req *http.Request) (Out, error) {
		// currently, the only way to genereate a zero value of a generic type is
		// like so:
		var zero Out

		userIDStr := req.Header.Get("X-User-Id")
		if userIDStr == "" {
			return zero, errors.New("please set the X-User-Id header with your ID")
		}

		userID, err := UserIDFromString(userIDStr)
		if err != nil {
			return zero, err
		}

		myAppCtx := MyAppContext{
			Context: ctx,
			UserID:  userID,
			Logger:  logrus.StandardLogger().WithField("context_id", uuid.NewV4().String()),
		}
		return h(myAppCtx, in, req)
	}
}

