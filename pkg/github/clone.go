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

type Cloner interface {
	Clone(context.Context, string, github.Repository)
}

type CloneConfig struct {
	URL  string
	Auth transport.AuthMethod
}

func NewCloneConfig(repo *github.Repository) CloneConfig {
	return CloneConfig{
		URL: repo.GetCloneURL(),
	}
}

func (cf CloneConfig) WithTokenAuth(token string) CloneConfig {
	cf.Auth = &http.BasicAuth{
		Username: "null",
		Password: token,
	}
	return cf
}

func (cf CloneConfig) WithBasicAuth(username, password string) CloneConfig {
	cf.Auth = &http.BasicAuth{
		Username: username,
		Password: password,
	}
	return cf
}

func (cf CloneConfig) WithAuth(auth transport.AuthMethod) CloneConfig {
	cf.Auth = auth
	return cf
}

type ErrorWithMeta struct {
	Error error

	ErrorMeta
}

type ErrorMeta struct {
	RepositoryName string
}

// CloneRepositories clones repositories as quickly as possible as subdirectories
// where the subdirectory name is the repository name, with the parent directory
// specified by dir. CloneRepositories guarantees that the repositories will be
// cloned as long as an error is not incurred.
func CloneRepositories(ctx context.Context, dir string, repos []*github.Repository, errCh chan ErrorWithMeta, cfg CloneConfig) {
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
