package github

import (
	"context"
	"fmt"

	"github.com/mccurdyc/neighbor/sdk/project"
	"github.com/mccurdyc/neighbor/sdk/search"
)

type searchMethodEntity = string

const (
	commit      searchMethodEntity = "commit"
	issue       searchMethodEntity = "issue"
	pullRequest searchMethodEntity = "pull_request"
	topic       searchMethodEntity = "topic"
	textMatch   searchMethodEntity = "text_match"
	label       searchMethodEntity = "label"
)

func Factory(ctx context.Context, conf *search.BackendConfig) (search.Backend, error) {
	if len(conf.AuthMethod) == 0 {
		// auth method required for GitHub code search - https://developer.github.com/v3/search/#search-code
		if conf.SearchMethod == search.CodeSearch {
			return nil, fmt.Errorf("auth method required for code search")
		}
	}

	var entity searchMethodEntity
	if conf.SearchMethod == search.VersionSearch {
		if conf.Config["version_entity"] == "" {
			return nil, fmt.Errorf("version_entity required with VersionSearch search method")
		}

		entity = conf.Config["version_entity"]
	}

	if conf.SearchMethod == search.MetaSearch {
		if conf.Config["meta_entity"] == "" {
			return nil, fmt.Errorf("meta_entity required with MetaSearch search method")
		}

		entity = conf.Config["meta_entity"]
	}

	return &Backend{
		searchMethod:       conf.SearchMethod,
		searchMethodEntity: entity,
	}, nil
}

type Backend struct {
	searchMethod       search.SearchMethod
	searchMethodEntity searchMethodEntity
}

func (b *Backend) Search(ctx context.Context, query string, numDesiredResults int) ([]project.Backend, error) {
	switch b.searchMethod {
	case search.ProjectSearch:
		searchRepository()
	case search.CodeSearch:
		searchCode()
	case search.MetaSearch:
	}

	return nil, nil // TODO: fix
}
