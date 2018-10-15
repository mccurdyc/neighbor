package external

import (
	// stdlib

	"os"
	"os/exec"

	// external

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

			// I do have concerns setting an environment variable if we do end up processing
			// projects concurrently. This environment variable could be overwritten in
			// another goroutine.
			// err = os.Setenv("COVERPROFILE_OUT_PATH", p.Directory)
			//
			// ctx.Logger.Infof("setting COVERPROFILE_OUT_PATH for %s to (%s)", p.Name, p.Directory)
			// if err != nil {
			// 	ctx.Logger.Error(err)
			// 	continue
			// }

			// we can't parse the command outside of this loop because exec.Command creates
			// a pointer to a Cmd and if you call Run() on that command, it will say
			// that it is already processing.
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
