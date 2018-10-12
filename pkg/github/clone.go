package github

import (
	// stdlib

	"fmt"
	"os"

	// external
	"github.com/google/go-github/github"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"

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

// CloneFromResult creates temporary directories where the base path is that of os.TempDir
// and the rest of the path is the Name of the repository. After creating a temporary
// directory, a project is cloned into that directory. After creating temp directories
// and cloning projects into the respective directory, the context is updated
// with the project names and the temporary directories.
func CloneFromResult(ctx *neighbor.Ctx, c *github.Client, d interface{}) {
	switch t := d.(type) {
	case *github.RepositoriesSearchResult:
		for _, r := range t.Repositories {

			dir := fmt.Sprintf("%s/%s", ctx.ExtResultDir, *r.Name)
			ctx.Logger.Infof("created directory: %s", dir)

			sshAuth, err := ssh.NewSSHAgentAuth("git") // username has to be "git" - https://github.com/src-d/go-git/issues/637
			if err != nil {
				ctx.Logger.Errorf("failed to create new ssh agent with error %+v", err)
				return
			}

			_, err = git.PlainClone(dir, false, &git.CloneOptions{
				Auth:     sshAuth,
				URL:      r.GetSSHURL(),
				Progress: os.Stdout,
			})
			if err != nil {
				ctx.Logger.Errorf("failed to clone project %s with error %+v", *r.Name, err)
				return
			}

			ctx.Logger.Infof("cloned: %s", r.GetCloneURL())

			ctx.ProjectDirMap[*r.Name] = dir
		}

		return
	case *github.CodeSearchResult:
		// needs implemented
	default:
		return
	}

	return
}
