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

	if conf.RetrievalFunc == nil {
		return nil, fmt.Errorf("retrieval function cannot be empty")
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

func (b *Backend) Name() string {
	return b.name
}

func (b *Backend) Version() string {
	return b.version
}

func (b *Backend) RetrievalFunc() retrieval.Backend {
	return b.retrievalFunc
}

func (b *Backend) SourceLocation() string {
	return b.SourceLocation()
}

func (b *Backend) LocalLocation() string {
	return b.localLocation
}

func (b *Backend) SetLocalLocation(l string) {
	b.localLocation = l
	return
}
