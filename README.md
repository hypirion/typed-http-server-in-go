# typed-http-server-in-go

This is a prototype of a typed HTTP server in Go via generics. Generics are
scheduled to come out with Go 1.18, and you can play around with this example
by using [gotip](https://pkg.go.dev/golang.org/dl/gotip).

I've written an entire blog post about my experimentation with this over here:
[Type-safe HTTP Servers in Go via
Generics](https://hypirion.com/musings/type-safe-http-servers-in-go-via-generics).
The important part is that I think values you attach to a context should be
typed, and this seemed to be the most reasonable approach to do so.

All middleware, encoding and decoding is done with type safe transforms that
doesn't use any kind of reflection. The end result is that you can make an HTTP
server that has HTTP endpoints looking like this:

```go
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
```

With as little middleware as this:

```go
func handlerWithMiddleware[In, Out any](handler typedhttp.Handler[middleware.MyAppContext,In,Out]) http.Handler {
	withAppContext := middleware.AttachUserAndLogger[context.Context,In,Out](handler)
	return httpjson.HandleTyped(withAppContext)
}
```


## Running

Assuming you have `gotip` installed, you can run the entire thing by going to
the directory `cmd/server`, then run `gotip run server.go`.

If you then pop over to a different terminal, you can run a call to the server
like so:

```shell
curl localhost:8080/ping -XPOST --header "X-User-Id: 123" --data-binary '{"message": "Hello"}'
```

## Structure

I've tried to group code into different packages, all of which resides in the
directory `pkg`. I suspect most of it will be code local to your project, but
I've added comments to all the files to denote whether I think it will be a
project-specific thing or a library you can import.

But who knows, this is only a prototype, and whether we'll even use this format
at all is not something I'm sure of.

## License

Typically I'd say "Copyright © 2021 Jean Niklas L'orange", but I've waived those
rights by applying a CC0 license to this repo. Do whatever you want with the
code that resides here, although I don't mind a referral back to this repo or to
my original blog post.
