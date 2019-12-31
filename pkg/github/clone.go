package github

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/google/go-github/github"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

const numWorkers = 5

// Cloner is the interface required to clone a repository.
type Cloner interface {
	Clone(context.Context, string, github.Repository)
}

// CloneConfig is used to configure cloning.
type CloneConfig struct {
	url  string
	auth transport.AuthMethod
}

// NewCloneConfig creates a new, empty, CloneConfig.
func NewCloneConfig() CloneConfig {
	return CloneConfig{}
}

// WithRepoURL returns a new CloneConfig with the url field set to the repository's
// CloneURL.
func (cf CloneConfig) WithRepoURL(repo *github.Repository) CloneConfig {
	cf.url = repo.GetCloneURL()
	return cf
}

// WithTokenAuth is http.BasicAuth configured so that an API token can be used
// instead of a username and password. This is how GitHub handles token authentication
// https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line#using-a-token-on-the-command-line
func (cf CloneConfig) WithTokenAuth(token string) CloneConfig {
	cf.auth = &http.BasicAuth{
		Username: "null",
		Password: token,
	}
	return cf
}

// WithBasicAuth is authentication via username and password.
func (cf CloneConfig) WithBasicAuth(username, password string) CloneConfig {
	cf.auth = &http.BasicAuth{
		Username: username,
		Password: password,
	}
	return cf
}

// WithAuth sets the auth method to an arbitrary auth method that satisfies the
// transport.AuthMethod interface, specified by the consumer.
func (cf CloneConfig) WithAuth(auth transport.AuthMethod) CloneConfig {
	cf.auth = auth
	return cf
}

// ErrorWithMeta contains the error, or nil, incurred while cloning the repository
// with additional metadata about the cloned repository.
type ErrorWithMeta struct {
	Error error

	Meta
}

// Meta contains information about the cloned repository.
type Meta struct {
	// RepositoryName is the full repository name i.e., 'username/reponame'.
	RepositoryName string
	// ClonedDir is the filepath to the cloned repository.
	ClonedDir string
}

// CloneRepositories clones repositories as quickly as possible as subdirectories
// where the subdirectory name is the repository name, with the parent directory
// specified by dir. CloneRepositories guarantees that the repositories will be
// cloned as long as an error is not incurred. Errors incurred cloning a repository
// will not prevent attempts at cloning additional repositories. doneCh will
// be populated with the error value, or nil otherwise for each repository.
func CloneRepositories(ctx context.Context, dir string, repos []*github.Repository, doneCh chan ErrorWithMeta, cfg CloneConfig) {
	ch := make(chan github.Repository)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for w := 0; w < numWorkers; w++ {
		go func() {
			cloneRepositories(ctx, dir, ch, errCh, cfg)
			wg.Done()
		}()
	}

	for _, repo := range repos {
		ch <- *repo
	}

	close(ch)
	wg.Wait()
}

func cloneRepositories(ctx context.Context, dir string, repoCh <-chan github.Repository, errCh chan<- ErrorWithMeta, cfg CloneConfig) {
	for repo := range repoCh {

		cfg.URL = repo.GetCloneURL()

		err := clone(ctx, filepath.Join(dir, repo.GetFullName()), repo, cfg)
		if err != nil {
			errCh <- ErrorWithMeta{Error: err, ErrorMeta: ErrorMeta{RepositoryName: repo.GetFullName()}}
		}
	}
}

// Clone clones a project to the directory specified by dir with the project directory
// name being the 'username/repository'. Note that repository is not a sub-directory.
func Clone(ctx context.Context, dir string, repo github.Repository, cfg CloneConfig) error {
	return clone(ctx, dir, repo, cfg)
}

// clone clones a project to the directory specified by dir with the project directory
// name being the 'username/repository'. Note that repository is not a sub-directory.
func clone(ctx context.Context, dir string, repo github.Repository, cfg CloneConfig) error {
	opts := git.CloneOptions{
		URL: cfg.URL,
	}

	if cfg.Auth != nil {
		opts.Auth = cfg.Auth
	}

	err := opts.Validate()
	if err != nil {
		return fmt.Errorf("CloneOptions are invalid: %w", err)
	}

	_, err = git.PlainClone(dir, false, &opts)
	if err != nil {
		return fmt.Errorf("failed to clone '%s' to '%s/': %w", repo.GetFullName(), dir, err)
	}
	return nil
}
