package main

import (
	"context"
	"fmt"
	"time"
	"net/http"

	"github.com/hypirion/typed-http-server-in-go/pkg/net/typedhttp"
	"github.com/hypirion/typed-http-server-in-go/pkg/httpjson"
	"github.com/hypirion/typed-http-server-in-go/pkg/middleware"
	"github.com/sirupsen/logrus"
)

// curl localhost:8080/ping -XPOST --header "X-User-Id: 123" --data-binary '{"message": "Hello"}'
func main() {
	logrus.Info("server starting at localhost:8080")
	http.Handle("/ping", handlerWithMiddleware(pingHandler))

	logrus.Fatal(http.ListenAndServe(":8080", nil))
}

func handlerWithMiddleware[In, Out any](handler typedhttp.Handler[middleware.MyAppContext,In,Out]) http.Handler {
	withTimeout := middleware.SetHandlerTimeout(20 * time.Second, handler)
	withAppContext := middleware.AttachUserAndLogger[context.Context,In,Out](withTimeout)
	return httpjson.HandleTyped(withAppContext)
}

type pingInput struct {
	Message string `json:"message"`
}

type pingOutput struct {
	Message string `json:"message"`
	Author middleware.UserID `json:"author"`
}

func pingHandler(ctx middleware.MyAppContext, in pingInput, req *http.Request) (*pingOutput, error) {
	ctx.Logger.WithField("user_id", ctx.UserID).Infof("Ping pong, message %s", in.Message)
	return &pingOutput{
		Message: fmt.Sprintf("PONG! Received %q", in.Message),
		Author: ctx.UserID,
	}, nil
}
