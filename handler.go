package slash

import (
	"errors"

	"golang.org/x/net/context"
)

var (
	ErrNoHandler    = errors.New("slash: no handler")
	ErrUnauthorized = errors.New("slash: invalid token")
)

// Handler represents a Command handler.
type Handler interface {
	// ServeCommand runs the command. The handler should return a string
	// that will be sent back to the user.
	ServeCommand(context.Context, Command) (text string, err error)
}

// HandlerFunc is a function that implements the Handler interface.
type HandlerFunc func(context.Context, Command) (string, error)

func (fn HandlerFunc) ServeCommand(ctx context.Context, command Command) (string, error) {
	return fn(ctx, command)
}

// Mux is a Handler implementation that routes commands to Handlers.
type Mux struct {
	routes map[string]Handler
}

// NewMux returns a new Mux instance.
func NewMux() *Mux {
	return &Mux{
		routes: make(map[string]Handler),
	}
}

// Handle adds a Handler to handle the given command.
//
// Example
//
//	m.Handle("/deploy", DeployHandler)
func (m *Mux) Handle(command string, handler Handler) {
	m.routes[command] = handler
}

// Handler returns the Handler that can handle the given slash command. If no
// handler matches, nil is returned.
func (m *Mux) Handler(command Command) Handler {
	h, ok := m.routes[command.Command]
	if !ok {
		return nil
	}
	return h
}

// ServeCommand attempts to find a Handler to serve the Command. If no handler
// is found, an error is returned.
func (m *Mux) ServeCommand(ctx context.Context, command Command) (string, error) {
	h := m.Handler(command)
	if h == nil {
		return "", ErrNoHandler
	}
	return h.ServeCommand(ctx, command)
}

// Authorize returns a new Handler that verifies that the token in the request
// matches the given token.
func Authorize(h Handler, token string) Handler {
	return HandlerFunc(func(ctx context.Context, command Command) (string, error) {
		if command.Token != token {
			return "", ErrUnauthorized
		}
		return h.ServeCommand(ctx, command)
	})
}
