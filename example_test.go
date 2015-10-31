package slash_test

import (
	"net/http"

	"golang.org/x/net/context"

	"github.com/remind101/slack-deployments/slash"
)

func Example() {
	m := slash.NewMux()
	m.Handle("/weather", slash.HandlerFunc(func(ctx context.Context, command slash.Command) (string, error) {
		return "cold!", nil
	}))

	s := slash.NewServer(slash.Authorize(m, "secret"))
	http.ListenAndServe(":8080", s)
}
