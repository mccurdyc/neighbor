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
		ctx.Logger.Infof("K: %s, V: %s\n", name, dir)

		err := os.Chdir(dir)
		if err != nil {
			ctx.Logger.Error(err)
			continue
		}

		err = ctx.TestCmd.Run()
		if err != nil {
			ctx.Logger.Error(err)
			continue
		}
	}
}
