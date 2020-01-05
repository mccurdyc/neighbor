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

func searchRepositories(ctx context.Context, c *github.Client, query string, numDesiredResults int) ([]project.Backend, *github.Response, error) {
	res := make([]project.Backend, 0, numDesiredResults)

	gitHubSearchOptions := github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: pageSize(numDesiredResults, maxPageSize),
		},
	}

	searchRes, resp, err := c.Search.Repositories(ctx, query, &gitHubSearchOptions)
	if err != nil {
		return res, resp, err
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

func getLatestCommit(ctx context.Context, c *github.Client, repo github.Repository) (*github.RepositoryCommit, error) {
	commits, _, err := c.Repositories.ListCommits(ctx, repo.GetOwner().String(), repo.GetName(), nil)
	if err != nil {
		return nil, err
	}

	if len(commits) < 1 {
		return nil, nil
	}

	return commits[0], nil
}

func searchCode(ctx context.Context) ([]project.Backend, *github.Response, error) {
	panic("not implemented")
}

func searchMeta(ctx context.Context, entity searchMethodEntity) ([]project.Backend, *github.Response, error) {
	switch entity {
	case topic:
		return searchTopic(ctx)
	case textMatch:
		return searchTextMatch(ctx)
	case label:
		return searchLabel(ctx)
	}

	return nil, nil, fmt.Errorf("search method entity unsupported")
}

func searchTopic(ctx context.Context) ([]project.Backend, *github.Response, error) {
	panic("not implemented")
}

func searchLabel(ctx context.Context) ([]project.Backend, *github.Response, error) {
	panic("not implemented")
}

func searchTextMatch(ctx context.Context) ([]project.Backend, *github.Response, error) {
	panic("not implemented")
}