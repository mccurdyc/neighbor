package github

import (
	// stdlib

	"fmt"
	"os"

	// external
	"github.com/google/go-github/github"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	// internal
	"github.com/mccurdyc/neighbor/pkg/neighbor"
)

// ClonedRepositoriesCtxKey is used to set a context key as it complains when you use a built-in
// type. This is what is suggested by bradfitz on GitHub https://github.com/golang/go/issues/17826#issuecomment-258946985.
// Error Message: "should not use basic type string as key in context.WithValue"
type ClonedRepositoriesCtxKey struct{}

// repoDirMap will store repository names as the key where the value will be the
// path to where the repository was cloned.
type repoDirMap map[string]string

// ExternalProject contains a GitHub project's name as where it was cloned to
type ExternalProject struct {
	Name      string
	Directory string
}

// CloneFromResult creates temporary directories where the base path is that of os.TempDir
// and the rest of the path is the Name of the repository. After creating a temporary
// directory, a project is cloned into that directory. After creating temp directories
// and cloning projects into the respective directory, the context is updated
// with the project names and the temporary directories.
func CloneFromResult(ctx *neighbor.Ctx, c *github.Client, d interface{}) <-chan ExternalProject {
	ch := make(chan ExternalProject)

	go func() {
		defer close(ch)

		switch t := d.(type) {
		case *github.RepositoriesSearchResult:
			for _, r := range t.Repositories {

				ctx.Logger.Debugf("%+v", r)
				dir := fmt.Sprintf("%s/%s", ctx.ExtResultDir, *r.Name)
				ctx.Logger.Infof("created directory: %s", dir)

				_, err := git.PlainClone(dir, false, &git.CloneOptions{
					// you must use BasicAuth with your GitHub Access Token as the password
					// and the Username can be anything.
					Auth: &http.BasicAuth{
						Username: "abc123", // yes, this can be anything except an empty string
						Password: ctx.GitHub.AccessToken,
					},
					URL:      r.GetCloneURL(),
					Progress: os.Stdout,
				})
				if err != nil {
					ctx.Logger.Errorf("failed to clone project %s with error: %+v", *r.Name, err)
					continue
				}

				ctx.Logger.Infof("cloned: %s", r.GetCloneURL())

				select {
				case ch <- ExternalProject{
					Name:      *r.Name,
					Directory: dir,
				}:
				case <-ctx.Context.Done():
					return
				}
			}
		}
	}()

	return ch
}
