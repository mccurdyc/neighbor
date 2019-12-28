package github

import (
	"context"

	"github.com/google/go-github/github"
)

func NewCodeSearcher(c *github.Client) *CodeSearcher {
	return &CodeSearcher{
		client: c,
	}
}

type CodeSearcher struct {
	client *github.Client
}

func (cs *CodeSearcher) search(ctx context.Context, query string, opts *github.SearchOptions) (Results, error) {
	res, resp, err := cs.client.Search.Code(ctx, query, opts)
	if err != nil {
		return Results{}, err
	}

	return cs.processResults(res, resp), nil
}

func (cs *CodeSearcher) processResults(r interface{}, resp *github.Response) Results {
	res := r.(*github.CodeSearchResult) // should panic if r cannot be type asserted

	repos := make([]*github.Repository, 0, res.GetTotal())

	for _, res := range res.CodeResults {
		repos = append(repos, res.GetRepository())
	}

	return Results{
		Repositories:     repos,
		response:         resp,
		rawSearchResults: r,
	}
}
