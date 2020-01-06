package github

import (
	"context"

	"github.com/google/go-github/github"
)

// Client is a minimal wrapper of the *github.Client.
// This makes it possible to test the GitHub search backend by mocking the GitHub client.
type Client struct {
	SearchService     SearchService
	RepositoryService RepositoryService

	*github.Client
}

// SearchService is the minimal search service interface required by the GitHub search backend.
// https://github.com/google/go-github/issues/113#issuecomment-454308733
type SearchService interface {
	// Repositories returns the repositories for a given search query.
	Repositories(context.Context, string, *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error)
}

// RepositoryService is the minimal repository search interface required by the GitHub search backend.
// https://github.com/google/go-github/issues/113#issuecomment-454308733
type RepositoryService interface {
	// ListCommits lists the commits for a specific repository.
	ListCommits(context.Context, string, string, *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error)
}
