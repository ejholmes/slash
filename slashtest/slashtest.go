// Package slashtest contains helpers for testing slash commands.
package slashtest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/ejholmes/slash"
)

// ResponseRecorder is a slash.Responder implementation for testing purposes. It
// records the responses in a channel that can then be received on to make
// assertions. It also attempts to mimick the behavior of Slack in that it will
// return an error if you try to send more than 5 responses.
type ResponseRecorder struct {
	Responses <-chan slash.Response

	// internal channel to send on.
	ch chan slash.Response
}

// NewRecorder returns a new ResponseRecorder with the Responses channel set to
// a buffered channel allowing 5 responses.
func NewRecorder() *ResponseRecorder {
	ch := make(chan slash.Response, slash.MaximumDelayedResponses)
	return &ResponseRecorder{
		Responses: ch,
		ch:        ch,
	}
}

// NewServer returns an httptest.Server that uses a ResponseRecorder as the
// http.Handler
func NewServer() (*ResponseRecorder, *httptest.Server) {
	r := NewRecorder()
	return r, httptest.NewServer(r)
}

// Respond sends the response on the Responses channel. If the channel is
// blocked, it returns an error.
func (r *ResponseRecorder) Respond(resp slash.Response) error {
	return r.add(resp)
}

// ServeHTTP makes the ResponseRecorder implement the http.Handler method so it
// can be used in combination with httptest.Server to record responses posted to
// Slack.
func (r *ResponseRecorder) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var resp slash.Response
	if err := json.NewDecoder(req.Body).Decode(&resp); err != nil {
		panic(err)
	}

	if err := r.add(resp); err != nil {
		panic(err)
	}
}

func (r *ResponseRecorder) add(resp slash.Response) error {
	select {
	case r.ch <- resp:
		return nil
	default:
		return fmt.Errorf("you can send a maximum of %d delayed responses", cap(r.ch))
	}
}
