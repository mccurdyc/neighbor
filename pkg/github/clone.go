package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

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

// Clone clones a project to the directory specified by dir with the project directory
// name being the 'username/repository'. Note that repository is not a sub-directory.
func Clone(ctx context.Context, dir string, repo github.Repository, cfg CloneConfig) error {
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
		return fmt.Errorf("failed to clone '%s' to '%s/': %w", *repo.Name, dir, err)
	}
	return nil
}
