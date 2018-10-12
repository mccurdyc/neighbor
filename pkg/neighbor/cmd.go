package neighbor

// stdlib
// external
// internal

// SetTestCmd sets the test command that will be run on external projects.
func SetTestCmd(c *Ctx) {
	c.TestCmd = c.Config.Contents.TestCmd
}
