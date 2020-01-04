package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/mccurdyc/neighbor/sdk/project"
	"github.com/mccurdyc/neighbor/sdk/search"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
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

	var auth transport.AuthMethod
	if strings.EqualFold(conf.AuthMethod, "basic") {
		username := conf.Config["username"]
		if len(username) == 0 {
			return nil, fmt.Errorf("username required for basic auth")
		}

		password := conf.Config["password"]
		if len(password) == 0 {
			return nil, fmt.Errorf("password required for basic auth")
		}

		auth = &http.BasicAuth{
			Username: username,
			Password: password,
		}
	}

	if strings.EqualFold(conf.AuthMethod, "token") {
		token := conf.Config["token"]

		if len(token) == 0 {
			return nil, fmt.Errorf("token required for token auth")
		}
		auth = &http.BasicAuth{
			Username: "null", // this can't be an empty string
			Password: token,
		}
	}

	return &Backend{
		auth:               auth,
		searchMethod:       conf.SearchMethod,
		searchMethodEntity: entity,
	}, nil
}

type Backend struct {
	auth               transport.AuthMethod
	searchMethod       search.SearchMethod
	searchMethodEntity searchMethodEntity
}

func (b *Backend) Search(ctx context.Context, query string, numDesiredResults int) ([]project.Backend, error) {
	// TODO: these search functions could possibly be an interface
	switch b.searchMethod {
	case search.ProjectSearch:
		searchRepository()
	case search.CodeSearch:
		searchCode()
	case search.MetaSearch:
	}

	return nil, nil // TODO: fix
}
