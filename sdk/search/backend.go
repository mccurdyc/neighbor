package search

import (
	"context"

	"github.com/mccurdyc/neighbor/sdk/project"
)

// SearchMethod is the method of search to be used to find projects.
type SearchMethod string

const (
	// ProjectSearch is a search method for explicitly searching for projects
	// (e.g., by name, etc.).
	ProjectSearch SearchMethod = "project"
	// CodeSearch is a search method for searching through code to identify projects
	// that meet the code search criteria (e.g., projects written in a language,
	// or that have a specific file or tests, etc.).
	CodeSearch SearchMethod = "code"
	// MetaSearch is a search method for searching through project meta information
	// (e.g., GitHub topics, GitHub pull requests, etc.).
	MetaSearch SearchMethod = "meta"
	// VersionSearch is a search method for searching through the revision history
	// of a project (e.g., Git commits, GitHub pull requests, etc.).
	VersionSearch SearchMethod = "version"
)

// Backend is the minimal interface for a search backend.
type Backend interface {
	Search(context.Context, string, int) ([]project.Backend, error)
}

// BackendConfig contains the configuration parameters for a search backend.
type BackendConfig struct {
	// AuthRequired is indicates the method, if one, to be used for authentication
	// by the retrieval backend.
	AuthMethod string

	// SearchMethod is the method of search to be used to find projects.
	SearchMethod SearchMethod

	// Config is for optional or secondary configuration.
	Config map[string]string
}

// Factory is a factory function for constructing a search backend.
type Factory func(context.Context, *BackendConfig) (Backend, error)
