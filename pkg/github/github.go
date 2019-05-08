package github

import (
	// stdlib
	"context"

	// external
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	// internal
)

// Connect returns an authenticated Github client.
func Connect(ctx context.Context, tkn string) *github.Client {
	if len(tkn) == 0 {
		return github.NewClient(nil)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: tkn},
	)

	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
