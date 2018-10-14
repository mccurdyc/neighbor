package external

import (
	// stdlib
	"fmt"
	"os"

	// external
	"github.com/google/uuid"

	// internal
	"github.com/mccurdyc/neighbor/pkg/neighbor"
)

// RunTests runs the tests of an external project using the context's TestCmd.
func RunTests(ctx *neighbor.Ctx) {
	for name, dir := range ctx.ProjectDirMap {
		ctx.Logger.Infof("running tests for %s", name)

		err := os.Chdir(dir)
		if err != nil {
			ctx.Logger.Errorf("error changing into %s directory: %+v", name, err)
			continue
		}

		// we need to append a globally unique identifier to the coverprofile
		// path because a project could have multiple coverage profiles from multiple
		// packages and we want to store them all in the root of the project with an
		// easily-identifiable name "neighbor-projectname-coverprofile-UUID.out"
		guuid, err := uuid.NewRandom()
		if err != nil {
			ctx.Logger.Errorf("error generating new random UUID: %+v", err)
		}

		cp := fmt.Sprintf("%s/neighbor-%s-coverprofile-%s.out", dir, name, guuid.String())

		err = os.Setenv("COVERPROFILE_OUT_PATH", cp)
		ctx.Logger.Infof("setting COVERPROFILE_OUT_PATH for %s to (%s)", name, cp)
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
