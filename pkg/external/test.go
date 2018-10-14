package external

import (
	// stdlib
	"fmt"
	"os"
	"os/exec"

	// external
	"github.com/google/uuid"

	// internal
	"github.com/mccurdyc/neighbor/pkg/github"
	"github.com/mccurdyc/neighbor/pkg/neighbor"
)

// RunTests runs the tests of an external project using the context's TestCmd.
func RunTests(ctx *neighbor.Ctx, ch <-chan github.ExternalProject) {
	run := func(ch <-chan github.ExternalProject) {
		for p := range ch {

			ctx.Logger.Infof("running tests for %s", p.Name)
			err := os.Chdir(p.Directory)
			if err != nil {
				ctx.Logger.Error(err)
				continue
			}

			if len(ctx.TestCmd) < 1 {
				ctx.Logger.Errorf("test command cannot be empty")
				return
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

			var cmd *exec.Cmd
			if len(ctx.TestCmd) == 1 {
				cmd = exec.Command(ctx.TestCmd[0])
			} else {
				cmd = exec.Command(ctx.TestCmd[0], ctx.TestCmd[1:]...)
			}

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				ctx.Logger.Errorf("failed to run test command with error %+v", err)
				continue
			}
		}

		select {
		case <-ctx.Context.Done():
			return
		}
	}

	go func() {
		run(ch)

		select {
		case <-ctx.Context.Done():
			return
		}
	}()
}
