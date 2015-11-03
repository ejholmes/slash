package slash

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"golang.org/x/net/context"
)

// Server adapts a Handler to be served over http.
type Server struct {
	Handler
}

// NewServer returns a new Server instance.
func NewServer(h Handler) *Server {
	return &Server{
		Handler: h,
	}
}

// ServeHTTP parses the Command from the incoming request then serves it using
// the Handler.
func (h *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	command, err := ParseRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responder := newResponder(command)
	if err := h.ServeCommand(context.Background(), responder, command); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Make future responses async.
	responder.async = true
	if resp := responder.reply; resp != nil {
		// If Respond was called before ServeCommand returned, this will
		// be used as the first response.
		json.NewEncoder(w).Encode(newResponse(*resp))
	}
}

// responder implements the Responder interface. It will hold a single Response
// object in memory while async is false. When async is set to true, it will
// fallback to the asyncResponder.
type responder struct {
	sync.Mutex
	async bool
	reply *Response

	// asyncResponder is the responder to use to send asynchronous
	// responses.
	asyncResponder Responder
}

func newResponder(command Command) *responder {
	return &responder{
		asyncResponder: NewAsyncResponder(command),
	}
}

func (r *responder) Respond(resp Response) error {
	r.Lock()
	defer r.Unlock()

	if r.async {
		return r.asyncResponder.Respond(resp)
	}

	if r.reply != nil {
		return errors.New("You can only reply once")
	}

	r.reply = &resp

	return nil
}

type response struct {
	ResponseType *string `json:"response_type,omitempty"`
	Text         string  `json:"text"`
}

func newResponse(resp Response) *response {
	r := &response{Text: resp.Text}
	if resp.InChannel {
		t := "in_channel"
		r.ResponseType = &t
	}
	return r
}
