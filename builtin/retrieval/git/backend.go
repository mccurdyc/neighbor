package git

import (
	"context"
	"fmt"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	"github.com/mccurdyc/neighbor/sdk/retrieval"
)

// Factory is the factory function for creating the backend for Git as a project
// retrieval method.
func Factory(ctx context.Context, conf *retrieval.BackendConfig) (retrieval.Backend, error) {
	var auth transport.AuthMethod

	if strings.EqualFold(conf.AuthMethod, "basic") {
		username := conf.Config["username"]
		if len(username) == 0 {
			return nil, fmt.Errorf("username required for basic auth")
		}

		password := conf.Config["password"]
		if len(password) == 0 {
			return nil, fmt.Errorf("password required for basic auth")
		}

		auth = &http.BasicAuth{
			Username: username,
			Password: password,
		}
	}

	if strings.EqualFold(conf.AuthMethod, "token") {
		token := conf.Config["token"]

		if len(token) == 0 {
			return nil, fmt.Errorf("token required for token auth")
		}
		auth = &http.BasicAuth{
			Username: "null", // this can't be an empty string
			Password: token,
		}
	}

	return &Backend{
		auth: auth,
	}, nil
}

// Backend is the backend for project retrieval using Git.
type Backend struct {
	auth transport.AuthMethod
}

// Retrieve clones a remote Git repository specified by src to a local dir.
func (b *Backend) Retrieve(ctx context.Context, src string, dir string) error {
	opts := git.CloneOptions{
		URL: src,
	}

	if b.auth != nil {
		opts.Auth = b.auth
	}

	err := opts.Validate()
	if err != nil {
		return err
	}

	_, err = git.PlainClone(dir, false, &opts)
	return err
}
