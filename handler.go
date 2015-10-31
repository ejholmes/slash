package slash

import (
	"errors"
	"regexp"

	"golang.org/x/net/context"
)

var (
	// ErrNoHandler is returned by Mux ServeCommand if a Handler isn't found
	// for the route.
	ErrNoHandler = errors.New("slash: no handler")

	// ErrUnauthorized is returned when the provided token in the request
	// does not match the expected secret.
	ErrUnauthorized = errors.New("slash: invalid token")
)

// Handler represents something that handles a slash command.
type Handler interface {
	// ServeCommand runs the command. The handler should return a string
	// that will be used as the reply to send back to the user, or an error.
	// If an error is returned, then the string value is what will be sent
	// to the user.
	ServeCommand(context.Context, Command) (reply string, err error)
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
//	m.Handle("/deploy", "token", DeployHandler)
func (m *Mux) Command(command, token string, handler Handler) {
	m.routes[command] = ValidateToken(handler, token)
}

// MatchText adds a route that matches when the text of the command matches the
// given regular expression. If the route matches and is called, slash.Matches
// will return the capture groups.
func (m *Mux) MatchText(re *regexp.Regexp, handler Handler) {
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

// ValidateToken returns a new Handler that verifies that the token in the
// request matches the given token.
func ValidateToken(h Handler, token string) Handler {
	return HandlerFunc(func(ctx context.Context, command Command) (string, error) {
		if command.Token != token {
			return "", ErrUnauthorized
		}
		return h.ServeCommand(ctx, command)
	})
}
