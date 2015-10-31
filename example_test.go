package slash_test

import (
	"net/http"
	"regexp"

	"github.com/ejholmes/slash"
	"golang.org/x/net/context"
)

func Example() {
	r := slash.NewMux()
	r.Command("/weather", slash.Authorize(Weather, "secret"))

	s := slash.NewServer(r)
	http.ListenAndServe(":8080", s)
}

// Weather is the primary slash handler for the /weather command.
func Weather(ctx context.Context, command slash.Command) (string, error) {
	h := slash.NewMux()

	var zipcodeRegex = regexp.MustCompile(`([0-9])`)
	h.MatchText(zipcodeRegex, slash.HandlerFunc(Zipcode))

	return h.ServeCommand(ctx, command)
}

// Zipcode is a slash handler that returns the weather for a zip code.
func Zipcode(ctx context.Context, command slash.Command) (string, error) {
	zip := slash.Matches(ctx)[0]
	return zip, nil
}
