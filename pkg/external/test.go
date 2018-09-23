package external

import (
	// stdlib
	"os"

	// external

	// internal
	"github.com/mccurdyc/neighbor/pkg/neighbor"
)

// RunTests runs the tests of an external project using the context's TestCmd.
func RunTests(ctx *neighbor.Ctx) {
	for name, dir := range ctx.ProjectDirMap {
		ctx.Logger.Infof("running tests for %s", name)
		err := os.Chdir(dir)
		if err != nil {
			ctx.Logger.Error(err)
			continue
		}

		ctx.TestCmd.Stdout = os.Stdout
		ctx.TestCmd.Stderr = os.Stderr

		if err := ctx.TestCmd.Run(); err != nil {
			ctx.Logger.Errorf("failed to run test command with error %+v", err)
			continue
		}
	}
}
