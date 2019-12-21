package github

import (
	"context"
	"math"

	"github.com/google/go-github/github"
)

const maxPageSize = 100 // https://developer.github.com/v3/#pagination

// SearchOptions specifies optional parameters to a GitHub search request.
type SearchOptions struct {
	NumDesiredResults int

	gitHubSearchOptions *github.SearchOptions
}

type SearchType string

const (
	Repository SearchType = "repository"
	Code       SearchType = "code"
)

type Result struct {
	Repository *github.Repository
}

type Searcher interface {
	search(context.Context, string, *github.SearchOptions) ([]Result, error)
}

func NewSearcher(c *github.Client, t SearchType) Searcher {
	switch t {
	case Repository:
		return NewRepositorySearcher(c)
	case Code:
		return NewCodeSearcher(c)
	}

	return nil
}

func Search(ctx context.Context, s Searcher, query string, opts *SearchOptions) ([]Result, error) {
	res := make([]Result, 0, opts.NumDesiredResults)

	var page int
	searchOpts := opts.gitHubSearchOptions
	searchOpts.PerPage = numPages(opts.NumDesiredResults, pageSize(opts.NumDesiredResults))

	for {
		searchRes, err := s.search(ctx, query, searchOpts)
		if err != nil {
			return res, err
		}

		for _, r := range searchRes {
			if len(res) >= opts.NumDesiredResults {
				return res, nil
			}

			res = append(res, r)
		}

		searchOpts.Page = page + 1
	}
}

func pageSize(desired int) int {
	if desired < maxPageSize {
		return desired
	}

	return maxPageSize
}

func numPages(desired, pageSize int) int {
	res := 1

	if desired > pageSize {
		res = int(math.Ceil(float64(desired) / float64(pageSize)))
	}

	return res
}

func NewCodeSearcher(c *github.Client) *CodeSearcher {
	return &CodeSearcher{
		client: c,
	}
}

type CodeSearcher struct {
	client *github.Client
}

func (c *CodeSearcher) search(ctx context.Context, query string, opts *github.SearchOptions) ([]Result, error) {
	res := make([]Result, 0)

	result, _, err := c.client.Search.Code(ctx, query, opts)
	if err != nil {
		return res, err
	}

	for _, r := range result.CodeResults {
		res = append(res, Result{
			Repository: r.Repository,
		})
	}

	return res, nil
}

func NewRepositorySearcher(c *github.Client) *RepositorySearcher {
	return &RepositorySearcher{
		client: c,
	}
}

type RepositorySearcher struct {
	client *github.Client
}

func (c *RepositorySearcher) search(ctx context.Context, query string, opts *github.SearchOptions) ([]Result, error) {
	res := make([]Result, 0)

	result, _, err := c.client.Search.Repositories(ctx, query, opts)
	if err != nil {
		return res, err
	}

	for _, r := range result.Repositories {
		res = append(res, Result{
			Repository: &r,
		})
	}

	return res, nil
}
