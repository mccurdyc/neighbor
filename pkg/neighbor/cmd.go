package neighbor

import (
	"strings"
)

// stdlib

// external
// internal

// SetExternalCmd sets the command that will be run on external projects.
func (c *Ctx) SetExternalCmd(s string) error {
	cmd := parseCmd(s)

	c.ExternalCmd = cmd
	return nil
}

func parseCmd(s string) []string {
	return strings.Split(s, " ")
}
