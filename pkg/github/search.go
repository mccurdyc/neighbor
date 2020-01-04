package github

import (
	"context"

	"github.com/google/go-github/github"
)

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
