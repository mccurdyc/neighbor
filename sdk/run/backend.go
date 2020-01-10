package run

import (
	"context"
	"io"
)

// Backend is the minimal interface for a search backend.
type Backend interface {
	Run() error
}

// BackendConfig contains the configuration parameters for a search backend.
type BackendConfig struct {
	// Name is the name of the command to be run.
	Name string
	// Args are the arguments for the command, specified by Name.
	Args string
	// Dir is the directory that should be used as the working directory for commands.
	Dir string
	// Stdout is where output should be written to.
	Stdout io.Writer
	// Stderr is where error output should be written to.
	Stderr io.Writer
	// Config is for optional or secondary configuration.
	Config map[string]string
}

// Factory is a factory function for constructing a search backend.
type Factory func(context.Context, *BackendConfig) (Backend, error)
