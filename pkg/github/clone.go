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
	Clone(context.Context, string, CloneConfig) error
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
func CloneRepositories(ctx context.Context, dir string, repos []*github.Repository, cloner Cloner, cfg CloneConfig) chan ErrorWithMeta {
	ch := make(chan github.Repository)
	doneCh := make(chan ErrorWithMeta)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	go func() {
		for _, repo := range repos {
			ch <- *repo
		}

		close(ch)
		// wait until all workers are done before closing the doneCh to avoid
		// panic due to sends on a closed channel.
		wg.Wait()
		close(doneCh)
	}()

	// start workers
	for w := 0; w < numWorkers; w++ {
		go func() {
			defer wg.Done()

			for repo := range ch {
				cfg.url = repo.GetCloneURL()
				dir := filepath.Join(dir, repo.GetFullName())

				err := cloner.Clone(ctx, dir, cfg)
				doneCh <- ErrorWithMeta{Error: err, Meta: Meta{RepositoryName: repo.GetFullName(), ClonedDir: dir}}
			}
		}()
	}

	return doneCh
}

type PlainCloner struct{}

func (pc *PlainCloner) Clone(ctx context.Context, dir string, cfg CloneConfig) error {
	opts := git.CloneOptions{
		URL: cfg.url,
	}

	if cfg.auth != nil {
		opts.Auth = cfg.auth
	}

	err := opts.Validate()
	if err != nil {
		return fmt.Errorf("CloneOptions are invalid: %w", err)
	}

	_, err = git.PlainClone(dir, false, &opts)
	return err
}
