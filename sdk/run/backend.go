package run

import (
	"context"
	"io"
)

// Backend is the minimal interface for a search backend.
type Backend interface {
	Run(context.Context, string) error
}

// BackendConfig contains the configuration parameters for a search backend.
type BackendConfig struct {
	// Cmd is the command to be run.
	Cmd string
	// Stdout is where output should be written to.
	Stdout io.Writer
	// Stderr is where error output should be written to.
	Stderr io.Writer
	// Config is for optional or secondary configuration.
	Config map[string]string
}

// Factory is a factory function for constructing a search backend.
type Factory func(context.Context, *BackendConfig) (Backend, error)
