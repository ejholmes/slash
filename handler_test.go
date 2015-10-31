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

	ctx := context.Background()
	h.On("ServeCommand", ctx, cmd).Return("", nil)

	_, err := m.ServeCommand(ctx, cmd)
	assert.NoError(t, err)

	h.AssertExpectations(t)
}

func TestMux_Command_NotFound(t *testing.T) {
	m := NewMux()

	cmd := Command{
		Command: "/deploy",
	}

	ctx := context.Background()
	_, err := m.ServeCommand(ctx, cmd)
	assert.Equal(t, err, ErrNoHandler)
}

func TestMux_MatchText_Found(t *testing.T) {
	h := new(mockHandler)
	m := NewMux()
	m.MatchText(regexp.MustCompile(`(?P<repo>\S+?) to (?P<environment>\S+?)$`), h)

	cmd := Command{
		Text: "acme-inc to staging",
	}

	ctx := context.Background()
	h.On("ServeCommand", ctx, cmd).Return("", nil)

	_, err := m.ServeCommand(ctx, cmd)
	assert.NoError(t, err)

	h.AssertExpectations(t)
}

func TestValidateToken(t *testing.T) {
	h := new(mockHandler)
	a := ValidateToken(h, "foo")

	ctx := context.Background()
	_, err := a.ServeCommand(ctx, Command{})
	assert.Equal(t, ErrUnauthorized, err)

	cmd := Command{
		Token: "foo",
	}
	h.On("ServeCommand", ctx, cmd).Return("", nil)
	_, err = a.ServeCommand(ctx, cmd)
	assert.NoError(t, err)
	h.AssertExpectations(t)
}
