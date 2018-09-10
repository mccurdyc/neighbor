package github

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Connect returns an authenticated Github client.
func Connect(ctx context.Context, tkn string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: tkn},
	)

	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
