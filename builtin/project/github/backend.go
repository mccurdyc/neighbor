package github

import (
	"context"
	"fmt"

	"github.com/mccurdyc/neighbor/sdk/project"
	"github.com/mccurdyc/neighbor/sdk/retrieval"
)

// Factory is a function for creating GitHub projects.
func Factory(ctx context.Context, conf *project.BackendConfig) (project.Backend, error) {
	if len(conf.Name) == 0 {
		return nil, fmt.Errorf("name cannot be empty")
	}

	if len(conf.SourceLocation) == 0 {
		return nil, fmt.Errorf("source location cannot be empty")
	}

	return &Backend{
		name:           conf.Name,
		version:        conf.Version,
		sourceLocation: conf.SourceLocation,
		retrievalFunc:  conf.RetrievalFunc,
	}, nil
}

// Backend is a GitHub project backend.
type Backend struct {
	name           string
	retrievalFunc  retrieval.Backend
	version        string
	sourceLocation string
	localLocation  string
}

// Name returns the name associated with a GitHub project.
func (b *Backend) Name() string {
	return b.name
}

// Version returns the version of the observed GitHub project (e.g., commit hash, semantic version, etc.).
func (b *Backend) Version() string {
	return b.version
}

// RetrievalFunc is the retrieval function that should be used to retrieve the project from GitHub.
// An example retrieval function could be Git.
func (b *Backend) RetrievalFunc() retrieval.Backend {
	return b.retrievalFunc
}

// SourceLocation is the source location, i.e., where the project was discovered (e.g., GitHub).
func (b *Backend) SourceLocation() string {
	return b.sourceLocation
}

// LocalLocation is the location on disk or where the project can be found in order
// to perform and evaluation or analysis of the project.
func (b *Backend) LocalLocation() string {
	return b.localLocation
}

// SetLocalLocation sets the local or on-disk location.
func (b *Backend) SetLocalLocation(l string) project.Backend {
	return &Backend{
		name:           b.Name(),
		retrievalFunc:  b.RetrievalFunc(),
		version:        b.Version(),
		sourceLocation: b.SourceLocation(),
		localLocation:  l,
	}
}
