package neighbor

import (
	// stdlib
	"fmt"
	"io"
	"os"
	// external
	// internal
)

// SwitchGoCmd temporarily switches from using the user's `go`, specified by the GOROOT environment
// variable to using `./bin/go-cover` which has the `-cover` and `-coverprofile`
// flags always enabled during testing.
//
// GOROOT is only required to be set when a user installs go in a custom location.
// https://stackoverflow.com/questions/7970390/what-should-be-the-values-of-gopath-and-goroot/18648321
func SwitchGoCmd(c *Ctx) error {
	// backup old go binary to temporary go.bak
	err := copyFile(fmt.Sprintf("%s/bin/go", c.GOROOT), fmt.Sprintf("%s/bin/go.bak", c.GOROOT))
	// copy the go binary with cover flags enabled to system go
	err = copyFile(fmt.Sprintf("%s/bin/go-cover", c.NeighborDir), fmt.Sprintf("%s/bin/go", c.GOROOT))
	return err
}

// CleanupGoCmd return the go command back to what it was before switching to use the
// go binary with the cover flags enabled.
func CleanupGoCmd(c *Ctx) error {
	err := copyFile(fmt.Sprintf("%s/bin/go.bak", c.GOROOT), fmt.Sprintf("%s/bin/go", c.GOROOT))
	return err
}

// SetTestCmd sets the test command that will be run on external projects.
func SetTestCmd(c *Ctx) {
	c.TestCmd = c.Config.Contents.TestCmd
}

func copyFile(src, dest string) error {
	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()

	destF, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destF.Close()

	_, err = io.Copy(destF, srcF)
	if err != nil {
		return err
	}

	err = destF.Sync()
	if err != nil {
		return err
	}
	return nil
}
