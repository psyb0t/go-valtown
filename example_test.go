package valtown_test

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-valtown"
)

// This is the README usage example. It talks to the live API, so it has no
// Output comment and only has to compile.
//
//nolint:testableexamples // needs a real API token, compile check only
func Example() {
	auth := valtown.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set(
				"Authorization",
				"Bearer "+os.Getenv("VALTOWN_API_TOKEN"),
			)

			return nil
		},
	)

	client, err := valtown.NewClientWithResponses(
		"https://api.val.town",
		auth,
	)
	if err != nil {
		panic(err)
	}

	res, err := client.MeGetWithResponse(context.Background())
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
