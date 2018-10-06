package neighbor

import (
	// stdlib
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	// external
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	// internal
	"github.com/mccurdyc/neighbor/pkg/config"
)

// Ctx is an object that contains information that will be used throughout the
// life of the neighor command. The idea was taken from the dep tool (https://github.com/golang/dep/blob/master/context.go#L23).
// This does NOT satisfice the context.Context interface (https://golang.org/pkg/context/#Context),
// therefore, it cannot be used as a context for methods or functions requiring a context.Context.
type Ctx struct {
	Config        *config.Config    // the query config created by the user
	Context       context.Context   // a context object required by the GitHub SDK
	GoRoot        string            // where the user's `go` binary is located
	Logger        *log.Logger       // the logger to be used throughout the project
	ProjectDirMap map[string]string // key: project name, value: absolute path to directory
	TestCmd       *exec.Cmd         // external project test command
}

// NewCtx creates a pointer to a new neighbor context.
func NewCtx() *Ctx {
	return &Ctx{}
}

// SwitchGoCmd temporarily switches from using the user's `go`, specified by the GOROOT environment
// variable to using `./bin/go-cover` which has the `-cover` and `-coverprofile`
// flags always enabled during testing.
//
// GOROOT is only required to be set when a user installs go in a custom location.
// https://stackoverflow.com/questions/7970390/what-should-be-the-values-of-gopath-and-goroot/18648321
//
// Thought: maybe we don't have to switch the binary, but instead just temporarily
// change GOROOT.
func SwitchGoCmd(c *Ctx) error {
	goRoot := "/usr/local/go"
	if p, ok := os.LookupEnv("GOROOT"); ok {
		goRoot = p
	}

	wd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "error getting neighbor directory")
	}

	// backup old go binary to temporary go.bak
	err = copyFile(fmt.Sprintf("%s/bin/go", goRoot), fmt.Sprintf("%s/bin/go.bak", goRoot))
	// copy the go binary with cover flags enabled to system go
	err = copyFile(fmt.Sprintf("%s/go-cover", wd), fmt.Sprintf("%s/bin/go", goRoot))
	return err
}

// CleanupGoCmd return the go command back to what it was before switching to use the
// go binary with the cover flags enabled.
func CleanupGoCmd(c *Ctx) error {
	goRoot := "/usr/local/go"
	if p, ok := os.LookupEnv("GOROOT"); ok {
		goRoot = p
	}

	err := copyFile(fmt.Sprintf("%s/bin/go.bak", goRoot), fmt.Sprintf("%s/bin/go", goRoot))
	return err
}

// SetTestCmd sets the test command that will be run on external projects.
func SetTestCmd(c *Ctx) {
	c.TestCmd = c.Config.Contents.TestCmd
}

func copyFile(src, dest string) error {
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}

	return nil
}
