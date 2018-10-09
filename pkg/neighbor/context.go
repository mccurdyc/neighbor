package neighbor

import (
	// stdlib
	"context"
	"os"
	"os/exec"

	// external

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
	GOROOT        string            // where the go binary and tools are located
	Logger        *log.Logger       // the logger to be used throughout the project
	NeighborDir   string            // the absolute path to neighbor project
	ProjectDirMap map[string]string // key: project name, value: absolute path to directory
	TestCmd       *exec.Cmd         // external project test command
}

// NewCtx creates a pointer to a new neighbor context.
func NewCtx() *Ctx {
	return &Ctx{}
}

// SetGOROOT sets the location of the go binary and tools. By default, it sets
// it to the default install location /usr/local/go, but if GOROOT is set, it
// will override the default.
func SetGOROOT(c *Ctx) {
	c.GOROOT = "/usr/local/go"
	if p := os.Getenv("GOROOT"); len(p) != 0 {
		c.GOROOT = p
	}
	return
}

// SetNeighborDir sets the neighbor directory field on the context to the absolute
// path of the neighbor project.
func SetNeighborDir(c *Ctx) error {
	wd, err := os.Getwd()

	c.NeighborDir = wd
	return err
}
