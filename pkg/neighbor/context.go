package neighbor

import (
	"os/exec"

	"github.com/mccurdyc/neighbor/pkg/config"
)

// Ctx is an object that contains information that will be used throughout the
// life of the neighor command. The idea was taken from the dep tool (https://github.com/golang/dep/blob/master/context.go#L23).
// This does NOT satisfice the context.Context interface (https://golang.org/pkg/context/#Context),
// therefore, it cannot be used as a context for methods or functions requiring a context.Context.
type Ctx struct {
	Config        config.Config     // the query config created by the user
	ProjectDirMap map[string]string // key: project name, value: absolute path to directory
	TestCmd       exec.Cmd          // external project test command
}

// NewCtx creates a pointer to a new neighbor context that contains the parsed config
// object.
func NewCtx(cfg config.Config) *Ctx {
	m := make(map[string]string)

	return &Ctx{
		Config:        cfg,
		ProjectDirMap: m,
	}
}

// AddToProjectDirMap populates the map with project names and the respective directory
// as the value.
func (c *Ctx) AddToProjectDirMap(k, v string) {
	c.ProjectDirMap[k] = v
}

// SetTestCmd sets the test command that will be run on external projects.
func (c *Ctx) SetTestCmd(cmd exec.Cmd) {
}
