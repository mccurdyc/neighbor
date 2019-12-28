package github

import (
	"context"
	"errors"

	"github.com/google/go-github/github"
)

// maxPageSize is the max number of results per page that GitHub returns.
// https://developer.github.com/v3/#pagination
const maxPageSize = 100

type searchOptions struct {
	numDesiredResults int
	maxPageSize       int

	gitHubSearchOptions github.SearchOptions
}

// SearchOptions returns the default searchOptions values (e.g., by default neighbor
// will attempt to return 100 repositories or one full page of results).
func SearchOptions() searchOptions {
	return searchOptions{
		numDesiredResults:   maxPageSize,
		maxPageSize:         maxPageSize,
		gitHubSearchOptions: github.SearchOptions{},
	}
}

// WithNumberOfResults sets the number of desired results that should be obtained.
func (so searchOptions) WithNumberOfResults(n int) searchOptions {
	so.numDesiredResults = n
	return so
}

// WithGitHubOptions sets the optional search parameters specified by GitHub.
// https://godoc.org/github.com/google/go-github/github#SearchOptions
func (so searchOptions) WithGitHubOptions(opts github.SearchOptions) searchOptions {
	so.gitHubSearchOptions = opts
	return so
}

// SearchType defines the supported search types.
// https://developer.github.com/v3/search/
//
// GitHub supports the following search types:
//  + Repository
//  + Code
//  + Commits
//  + Issues and Pull Request
//  + Users
//  + Topics
//  + Labels
//  + Text Match Metadata
type SearchType string

const (
	Repository SearchType = "repository"
	Code       SearchType = "code"
)

// Results contains the full list of repositories --- i.e., from multiple pages ---
// returned from a search, the response from GitHub which indicates how many total
// results and pagination information and the rawSearchResults which contains
// the search result information in addition to the repositories.
type Results struct {
	Repositories     []*github.Repository
	response         *github.Response
	rawSearchResults interface{}
}

// Searcher is an interface for searching and processing results.
type Searcher interface {
	search(context.Context, string, *github.SearchOptions) (Results, error)

	ResultProcessor
}

// ResultProcessor is an interface for processing the raw results from a search.
type ResultProcessor interface {
	processResults(interface{}, *github.Response) Results
}

// NewSearcher creates a new Searcher or indicates when a search type is not supported.
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

// ErrFewerResultsThanDesired is used to indicate that it was not possible to fulfill
// the request from the user (i.e., could not find the number of results specified
// by the user).
//
// This is important to specify because, for example, in research you might want
// to guarantee that you are analyzing _exactly_ the number of projects specifed
// or the search query may need to be tweaked.
var ErrFewerResultsThanDesired = errors.New("contains fewer results than desired")

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
			return res, ErrFewerResultsThanDesired
		}

		opts.gitHubSearchOptions.Page = page + 1
	}
}

// pageSize returns the minimal page size necessary to fulfill the request or the
// maximum page supported by GitHub.
// https://developer.github.com/v3/#pagination
func pageSize(desired, max int) int {
	if desired < max {
		return desired
	}

	return max
}
