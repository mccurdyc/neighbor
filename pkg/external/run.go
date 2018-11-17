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

const numWorkers = 3

// Run runs an arbitrary command specified in the Ctx on each project
// that is sent through the pipeline.
func Run(ctx *neighbor.Ctx, ch <-chan github.ExternalProject) chan struct{} {
	done := make(chan struct{})

	// wrap in a goroutine so that we can return the 'done' channel immediately
	go func() {
		for i := 0; i < numWorkers; i++ {
			// each worker should continue recieving projects until a quit signal is received
			go func() {
				for {
					select {
					case p := <-ch:
						run(ctx, p)
					case <-ctx.Context.Done():
						return
					}
				}
			}()
		}

		done <- struct{}{}
	}()

	return done
}

func run(ctx *neighbor.Ctx, p github.ExternalProject) {
	ctx.Logger.Infof("running external command on %s", p.Name)
	err := os.Chdir(p.Directory)
	if err != nil {
		ctx.Logger.Error(err)
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
	}
}
