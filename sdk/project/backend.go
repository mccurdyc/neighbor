package project

import (
	"context"

	"github.com/mccurdyc/neighbor/sdk/retrieval"
)

// Factory is a factory function for creating projects.
type Factory func(context.Context, *BackendConfig) (Backend, error)

// Backend is the minimal interface of a project.
type Backend interface {
	Name() string
	Version() string
	RetrievalFunc() retrieval.Backend
	SourceLocation() string
	LocalLocation() string
	SetLocalLocation(string) Backend
}

// BackendConfig is the configuration parameters for a project backend.
type BackendConfig struct {
	// Name is the name or an identifier of the project.
	Name string
	// Version is the current version of the project (e.g., semantic version, Git commit hash, etc.).
	Version string
	// SourceLocation is where (e.g., local file path, remote url, etc.) the project
	// can be found and retrieved from.
	SourceLocation string
	// RetrievalFunc is the function that was or can be used to retrieve a project.
	// Examples of retrieval funcitons could be git or a local copy.
	//
	// TODO: thought - Is the Retrival func for how they _retrieved_ the project
	// or, how someone _can retrieve_ the project?
	//
	// I think it is the latter because Search returns a project, but does not
	// retrieve the project.
	RetrievalFunc retrieval.Backend

	// Config is a way to set additional, optional and/or secondary configuration values.
	Config map[string]string
}
