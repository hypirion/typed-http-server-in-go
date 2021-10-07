// Package httpjson implements typedhttp.Handler for JSON payloads.
package httpjson

import (
	"context"
	"encoding/json"
	"net/http"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/hypirion/typed-http-server-in-go/pkg/net/typedhttp"
)

// I guess this code could either be local code or third party code, depending
// on how fine-grained control you want. People will probably stick to the
// defaults until they aren't enough I guess.

// HandleTyped takes a typed HTTP handler and returns a standard
// library-compatible handler that reads and writes JSON.
func HandleTyped[In, Out any](handler typedhttp.Handler[context.Context, In, Out]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		bs, err := ioutil.ReadAll(req.Body)
		if err != nil {
			writeError(w, err)
			return
		}
		var input In
		err = json.Unmarshal(bs, &input)
		if err != nil {
			writeError(w, err)
			return
		}
		out, err := handler(req.Context(), input, req)
		if err != nil {
			writeError(w, err)
			return
		}
		bs, err = json.Marshal(out)
		if err != nil {
			writeError(w, err)
			return
		}
		_, err = w.Write(bs)
		if err != nil {
			logrus.WithError(err).Warn("httpjson: could not write all the data back to the client")
		}
	})
}

type errorObject struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, err error) {
	logrus.WithError(err).Warn("error while handling response")
	w.WriteHeader(http.StatusInternalServerError)
	err = json.NewEncoder(w).Encode(errorObject{Error: err.Error()})
	if err != nil {
		logrus.WithError(err).Error("could not serialise json error object")
	}
}
