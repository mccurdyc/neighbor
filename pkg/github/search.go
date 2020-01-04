package github

import (
	"context"
	"errors"

	"github.com/google/go-github/github"
)
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
