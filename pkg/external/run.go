package external

import (
	// stdlib
	"os"
	"os/exec"
	"sync"

	// external
	// internal
	"github.com/mccurdyc/neighbor/pkg/github"
	"github.com/mccurdyc/neighbor/pkg/neighbor"
)

// Run runs an arbitrary command specified in the Ctx on each project
// that is sent through the pipeline.
func Run(ctx *neighbor.Ctx, ch <-chan github.ExternalProject) {
	run := func(ch <-chan github.ExternalProject) {
		for p := range ch {

			ctx.Logger.Infof("running external command on %s", p.Name)
			err := os.Chdir(p.Directory)
			if err != nil {
				ctx.Logger.Error(err)
				continue
			}

			if len(ctx.ExternalCmd) < 1 {
				ctx.Logger.Errorf("external command cannot be empty")
				return
			}

			// we can't parse the command outside of this loop because exec.Command creates
			// a pointer to a Cmd and if you call Run() on that command, it will say
			// that it is already processing.
			var cmd *exec.Cmd
			if len(ctx.ExternalCmd) == 1 {
				cmd = exec.Command(ctx.ExternalCmd[0])
			} else {
				cmd = exec.Command(ctx.ExternalCmd[0], ctx.ExternalCmd[1:]...)
			}

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				ctx.Logger.Errorf("failed to run external command with error %+v", err)
				continue
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		run(ch)
		wg.Done()

		select {
		case <-ctx.Context.Done():
			return
		}
	}()

	wg.Wait()
}
