package slash

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/net/context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServer_Reply(t *testing.T) {
	h := new(mockHandler)
	s := &Server{
		Handler: h,
	}

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", strings.NewReader(testForm))

	h.On("ServeCommand",
		context.Background(),
		mock.AnythingOfType("Command"),
	).Return(Reply("ok"), nil)

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, `{"text":"ok"}`+"\n", resp.Body.String())
}

func TestServer_Say(t *testing.T) {
	h := new(mockHandler)
	s := &Server{
		Handler: h,
	}

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", strings.NewReader(testForm))

	h.On("ServeCommand",
		context.Background(),
		mock.AnythingOfType("Command"),
	).Return(Say("ok"), nil)

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, `{"response_type":"in_channel","text":"ok"}`+"\n", resp.Body.String())
}

func TestServer_Err(t *testing.T) {
	h := new(mockHandler)
	s := &Server{
		Handler: h,
	}

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", strings.NewReader(testForm))

	errBoom := errors.New("boom")
	h.On("ServeCommand",
		context.Background(),
		mock.AnythingOfType("Command"),
	).Return(Reply(""), errBoom)

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
