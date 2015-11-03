package slash

import (
	"regexp"
	"testing"

	"golang.org/x/net/context"

	"github.com/stretchr/testify/assert"
)

func TestMux_Command_Found(t *testing.T) {
	h := new(mockHandler)
	m := NewMux()
	m.Command("/deploy", "token", h)

	cmd := Command{
		Token:   "token",
		Command: "/deploy",
	}

	resp := new(mockResponder)
	ctx := context.Background()
	h.On("ServeCommand",
		WithParams(ctx, make(map[string]string)),
		resp,
		cmd,
	).Return(nil)

	err := m.ServeCommand(ctx, resp, cmd)
	assert.NoError(t, err)

	h.AssertExpectations(t)
}

func TestMux_Command_NotFound(t *testing.T) {
	m := NewMux()

	cmd := Command{
		Command: "/deploy",
	}

	resp := new(mockResponder)
	ctx := context.Background()
	err := m.ServeCommand(ctx, resp, cmd)
	assert.Equal(t, err, ErrNoHandler)
}

func TestMux_MatchText_Found(t *testing.T) {
	h := new(mockHandler)
	m := NewMux()
	m.MatchText(regexp.MustCompile(`(?P<repo>\S+?) to (?P<environment>\S+?)$`), h)

	cmd := Command{
		Text: "acme-inc to staging",
	}

	resp := new(mockResponder)
	ctx := context.Background()
	h.On("ServeCommand",
		WithParams(ctx, map[string]string{"repo": "acme-inc", "environment": "staging"}),
		resp,
		cmd,
	).Return(nil)

	err := m.ServeCommand(ctx, resp, cmd)
	assert.NoError(t, err)

	h.AssertExpectations(t)
}

func TestValidateToken(t *testing.T) {
	h := new(mockHandler)
	a := ValidateToken(h, "foo")

	resp := new(mockResponder)
	ctx := context.Background()
	err := a.ServeCommand(ctx, resp, Command{})
	assert.Equal(t, ErrUnauthorized, err)

	cmd := Command{
		Token: "foo",
	}
	h.On("ServeCommand", ctx, resp, cmd).Return(nil)
	err = a.ServeCommand(ctx, resp, cmd)
	assert.NoError(t, err)
	h.AssertExpectations(t)
}

func TestMatchTextRegexp(t *testing.T) {
	re := regexp.MustCompile(`(?P<repo>\S+?) to (?P<environment>\S+?)(!)?$`)
	m := MatchTextRegexp(re)

	_, ok := m.Match(Command{Text: "foo"})
	assert.False(t, ok)

	params, ok := m.Match(Command{Text: "acme-inc to staging"})
	assert.True(t, ok)
	assert.Equal(t, map[string]string{"repo": "acme-inc", "environment": "staging"}, params)

	params, ok = m.Match(Command{Text: "acme-inc to staging!"})
	assert.True(t, ok)
	assert.Equal(t, map[string]string{"repo": "acme-inc", "environment": "staging"}, params)
}
