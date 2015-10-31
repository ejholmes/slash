package slash

import (
	"net/http"
	"net/url"

	"golang.org/x/net/context"
)

// Command represents an incoming Slash Command request.
type Command struct {
	Token string

	TeamID     string
	TeamDomain string

	ChannelID   string
	ChannelName string

	UserID   string
	UserName string

	Command string
	Text    string
}

// CommandFromValues returns a Command object from a url.Values object.
func CommandFromValues(v url.Values) Command {
	return Command{
		Token:       v.Get("token"),
		TeamID:      v.Get("team_id"),
		TeamDomain:  v.Get("team_domain"),
		ChannelID:   v.Get("channel_id"),
		ChannelName: v.Get("channel_name"),
		UserID:      v.Get("user_id"),
		UserName:    v.Get("user_name"),
		Command:     v.Get("command"),
		Text:        v.Get("text"),
	}
}

// ParseRequest parses the form an then returns the extracted Command.
func ParseRequest(r *http.Request) (Command, error) {
	err := r.ParseForm()
	return CommandFromValues(r.Form), err

}

// Params returns the match groups from a regular expression match.
func Params(ctx context.Context) map[string]string {
	params := make(map[string]string)
	return params
}

func WithParams(ctx context.Context, params map[string]string) context.Context {
	return ctx
}
