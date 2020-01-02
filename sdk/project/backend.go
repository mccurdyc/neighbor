package project

import "context"

// Backend is the minimal interface necessary to implement a project backend.
// Projects are the base entity that neighbor operates on.
// 	1. neighbor searches for projects
// 	2. neighbor retrieves projects
// 	3. neighbor analyzes, evaluates or executes a binary on projects.
type Backend interface {
	// GetName returns the name of the project.
	GetName() string
	// GetLocalLocation identifies where the project currently lives on the machine
	// running neighbor.
	GetLocalLocation() string
	// GetSourceLocation identifies where the project came from (e.g., Git clone URL,
	// another on-prem server, etc.).
	GetSourceLocation() string
	// GetVersion is means to identify the version of the project (e.g., commit hash
	// or semantic version for projects on GitHub).
	GetVersion() string
}

// BackendConfig contains configuration parameters used in the factory func to
// instantiate project backends.
type BackendConfig struct {
}

// Factory is the factory function to create a project backend.
type Factory func(context.Context, *BackendConfig) (Backend, error)
