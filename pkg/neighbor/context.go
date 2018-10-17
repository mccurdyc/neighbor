package neighbor

import (
	// stdlib
	"context"
	"os"

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
	Config       *config.Config  // the query config created by the user
	Context      context.Context // a context object required by the GitHub SDK
	GitHub       GitHubDetails
	Logger       *log.Logger // the logger to be used throughout the project
	NeighborDir  string
	ExtResultDir string   // where the external projects and test results will be stored
	TestCmd      []string // external project test command and args
}

// GitHubDetails are GitHub-specifc details necessary throughout the project
type GitHubDetails struct {
	AccessToken string
}

// NewCtx creates a pointer to a new neighbor context.
func NewCtx() *Ctx {
	return &Ctx{}
}

// CreateExternalResultDir creates the external projects and results directory if
// it doesn't exist.
func (ctx *Ctx) CreateExternalResultDir() error {
	_, err := os.Stat(ctx.ExtResultDir)
	if os.IsNotExist(err) {
		return os.Mkdir(ctx.ExtResultDir, os.ModePerm)
	}
	return nil
}
