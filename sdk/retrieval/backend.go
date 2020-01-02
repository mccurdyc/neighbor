package retrieval

import "context"

// Backend is the minimal interface that must be implemented by a retriever.
// Retrievers are responsible for obtaining projects, the base entity on which
// neighbor operates. An example retriever is Git.
type Backend interface {
	Retrieve(context.Context, string, string) error
}

// BackendConfig contains the configuration parameters for a retrieval backend.
type BackendConfig struct {
	// AuthRequired is indicates the method, if one, to be used for authentication
	// by the retrieval backend.
	AuthMethod string

	// Config is for optional or secondary configuration.
	Config map[string]string
}

// Factory is a factory function for constructing retrievers backend.
type Factory func(context.Context, *BackendConfig) (Backend, error)
