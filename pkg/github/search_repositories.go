package github

import (
	"context"

	"github.com/google/go-github/github"
)

func NewRepositorySearcher(c *github.Client) *RepositorySearcher {
	return &RepositorySearcher{
		client: c,
	}
}

type RepositorySearcher struct {
	client *github.Client
}

func (rs *RepositorySearcher) search(ctx context.Context, query string, opts *github.SearchOptions) (Results, error) {
	res, resp, err := rs.client.Search.Repositories(ctx, query, opts)
	if err != nil {
		return Results{}, err
	}

	return rs.processResults(res, resp), nil
}

func (rs *RepositorySearcher) processResults(r interface{}, resp *github.Response) Results {
	res := r.(*github.RepositoriesSearchResult) // should panic if r cannot be type asserted

	repos := make([]*github.Repository, 0, res.GetTotal())

	for i := range res.Repositories {
		repos = append(repos, &res.Repositories[i])
	}

	return Results{
		Repositories:     repos,
		response:         resp,
		rawSearchResults: r,
	}
}
