package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	githubProject "github.com/mccurdyc/neighbor/builtin/project/github"
	"github.com/mccurdyc/neighbor/sdk/project"
)

// maxPageSize is the max number of results per page that GitHub returns.
// https://developer.github.com/v3/#pagination
const maxPageSize = 100

type searchOptions struct {
	numDesiredResults int
	maxPageSize       int
}

func searchRepositories(ctx context.Context, c Client, query string, numDesiredResults int, opts *github.SearchOptions) ([]project.Backend, *github.Response, error) {
	res := make([]project.Backend, 0, numDesiredResults)

	searchRes, resp, err := c.SearchService.Repositories(ctx, query, opts)
	if err != nil {
		return res, resp, err
	}

	if searchRes == nil {
		return res, resp, fmt.Errorf("empty repository response")
	}

	for _, repo := range searchRes.Repositories {
		var version string
		latest, _ := getLatestCommit(ctx, c, repo)
		if latest != nil {
			version = latest.GetSHA()
		}

		p, err := githubProject.Factory(ctx, &project.BackendConfig{
			Name:           repo.GetFullName(),
			Version:        version,
			SourceLocation: repo.GetCloneURL(),
		})
		if err != nil {
			continue
		}

		res = append(res, p)
	}

	return res, resp, nil
}

func getLatestCommit(ctx context.Context, c Client, repo github.Repository) (*github.RepositoryCommit, error) {
	commits, _, err := c.RepositoryService.ListCommits(ctx, repo.GetOwner().GetName(), repo.GetName(), nil)
	if err != nil {
		return nil, err
	}

	if len(commits) < 1 {
		return nil, nil
	}

	return commits[0], nil
}

func searchCode(ctx context.Context, opts *github.SearchOptions) ([]project.Backend, *github.Response, error) {
	panic("not implemented")
}

func searchMeta(ctx context.Context, entity searchMethodEntity, opts *github.SearchOptions) ([]project.Backend, *github.Response, error) {
	switch entity {
	case topic:
		return searchTopic(ctx, opts)
	case textMatch:
		return searchTextMatch(ctx, opts)
	case label:
		return searchLabel(ctx, opts)
	}

	return nil, nil, fmt.Errorf("search method entity unsupported")
}

func searchTopic(ctx context.Context, opts *github.SearchOptions) ([]project.Backend, *github.Response, error) {
	panic("not implemented")
}

func searchLabel(ctx context.Context, opts *github.SearchOptions) ([]project.Backend, *github.Response, error) {
	panic("not implemented")
}

func searchTextMatch(ctx context.Context, opts *github.SearchOptions) ([]project.Backend, *github.Response, error) {
	panic("not implemented")
}
