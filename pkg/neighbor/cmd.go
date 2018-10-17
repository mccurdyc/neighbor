package neighbor

// stdlib
import "strings"

// external
// internal

// SetTestCmd sets the test command that will be run on external projects.
func SetTestCmd(c *Ctx) {
	c.TestCmd = strings.Split(c.Config.Contents.TestCmdStr, " ")
}
