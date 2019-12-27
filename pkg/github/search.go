package github

import (
	"context"
	"errors"

	"github.com/google/go-github/github"
)

const maxPageSize = 100 // https://developer.github.com/v3/#pagination

type searchOptions struct {
	numDesiredResults int
	maxPageSize       int

	gitHubSearchOptions github.SearchOptions
}

func SearchOptions() searchOptions {
	return searchOptions{
		numDesiredResults:   maxPageSize,
		maxPageSize:         maxPageSize,
		gitHubSearchOptions: github.SearchOptions{},
	}
}

func (so searchOptions) WithNumberOfResults(n int) searchOptions {
	so.numDesiredResults = n
	return so
}

func (so searchOptions) WithGitHubOptions(opts github.SearchOptions) searchOptions {
	so.gitHubSearchOptions = opts
	return so
}

type SearchType string

const (
	Repository SearchType = "repository"
	Code       SearchType = "code"
)

type Results struct {
	Repositories     []*github.Repository
	response         *github.Response
	rawSearchResults interface{}
}

type Searcher interface {
	search(context.Context, string, *github.SearchOptions) (Results, error)

	ResultProcessor
}

type ResultProcessor interface {
	processResults(interface{}, *github.Response) Results
}

func NewSearcher(c *github.Client, t SearchType) (Searcher, error) {
	switch t {
	case Repository:
		return NewRepositorySearcher(c), nil
	case Code:
		return NewCodeSearcher(c), nil
	default:
		return nil, errors.New("unsupported search type")
	}
}

var ErrRequestNotFulfilled = errors.New("contains fewer results than desired")

func Search(ctx context.Context, s Searcher, query string, opts searchOptions) ([]*github.Repository, error) {
	res := make([]*github.Repository, 0, opts.numDesiredResults)
	var page int

	opts.gitHubSearchOptions.PerPage = pageSize(opts.numDesiredResults, opts.maxPageSize)

	for {
		searchRes, err := s.search(ctx, query, &opts.gitHubSearchOptions)
		if err != nil {
			return res, err
		}

		for _, r := range searchRes.Repositories {
			if len(res) >= opts.numDesiredResults {
				return res, nil
			}

			res = append(res, r)
		}

		if searchRes.response == nil || searchRes.response.NextPage == 0 {
			return res, ErrRequestNotFulfilled
		}

		opts.gitHubSearchOptions.Page = page + 1
	}
}

func pageSize(desired, max int) int {
	if desired < max {
		return desired
	}

	return max
}
