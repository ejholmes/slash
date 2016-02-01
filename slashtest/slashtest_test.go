package slashtest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/ejholmes/slash"
	"github.com/stretchr/testify/assert"
)

func TestResponseServer(t *testing.T) {
	h := slash.NewServer(slash.HandlerFunc(func(ctx context.Context, r slash.Responder, c slash.Command) error {
		return r.Respond(slash.Reply("Hey"))
	}))

	// Responses from the above handler will be posted here.
	s := NewServer()
	defer s.Close()

	req, _ := http.NewRequest("POST", "/", strings.NewReader(fmt.Sprintf("response_url=%s", s.URL)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()

	h.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	select {
	case resp := <-s.Responses():
		assert.Equal(t, "Hey", resp.Text)
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}
