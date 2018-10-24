package neighbor

// stdlib
import "strings"

// external
// internal

// SetExternalCmd sets the command that will be run on external projects.
func SetExternalCmd(c *Ctx) {
	c.ExternalCmd = strings.Split(c.Config.Contents.ExternalCmdStr, " ")
}
