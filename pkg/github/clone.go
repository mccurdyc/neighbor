package github

import (
	// stdlib

	"context"
	"fmt"

	// external

	"github.com/google/go-github/github"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	// internal
)

// repoDirMap will store repository names as the key where the value will be the
// path to where the repository was cloned.
type repoDirMap map[string]string

// ExternalProject contains a GitHub project's name as where it was cloned to
type ExternalProject struct {
	Name      string
	Directory string
}

type Cloner interface {
	Clone(context.Context, string, github.Repository)
}

func WithCloneURL(opts *git.CloneOptions, url string) *git.CloneOptions {
	opts.URL = url
	return opts
}

func WithTokenAuth(opts *git.CloneOptions, token string) *git.CloneOptions {
	opts.Auth = withAuth(&http.BasicAuth{
		Username: "null",
		Password: token,
	})
	return opts
}

func WithBasicAuth(opts *git.CloneOptions, username, password string) *git.CloneOptions {
	opts.Auth = withAuth(&http.BasicAuth{
		Username: username,
		Password: password,
	})
	return opts
}

func WithAuth(opts *git.CloneOptions, auth transport.AuthMethod) *git.CloneOptions {
	opts.Auth = withAuth(auth)
	return opts
}

func withAuth(auth transport.AuthMethod) transport.AuthMethod {
	return auth
}

// Clone clones a project to the directory specified by dir with the project directory
// name being the 'username/repository'. Note that repository is not a sub-directory.
func Clone(ctx context.Context, dir string, repo github.Repository, opts *git.CloneOptions) error {
	err := opts.Validate()
	if err != nil {
		return fmt.Errorf("CloneOptions are invalid: %w", err)
	}

	_, err = git.PlainClone(dir, false, opts)
	if err != nil {
		return fmt.Errorf("failed to clone project %s to %s: %w", *repo.Name, dir, err)
	}
	return nil
}

// getCloneURL returns a GitHub git clone URL e.g., https://github.com/mccurdyc/neighbor.git
func getCloneURL(repo github.Repository) string {
	url := repo.GetCloneURL()
	if url == "" {
		url = fmt.Sprintf("%s.git", repo.GetHTMLURL())
	}
	return url
}
