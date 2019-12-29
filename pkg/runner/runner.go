package runner

import (
	"os"
)

type Runner interface {
	Run() error
}

// RunInDir runs an arbitrary Run function in the working directory specified by dir.
func RunInDir(dir string, r Runner) error {
	err := os.Chdir(dir)
	if err != nil {
		return err
	}

	return r.Run()
}
