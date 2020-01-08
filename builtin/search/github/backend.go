package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
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

// ErrFewerResultsThanDesired is used to indicate that it was not possible to fulfill
// the request from the user (i.e., could not find the number of results specified
// by the user).
//
// This is important to specify because, for example, in research you might want
// to guarantee that you are analyzing _exactly_ the number of projects specifed
// or the search query may need to be tweaked.
var ErrFewerResultsThanDesired = fmt.Errorf("contains fewer results than desired")

// Factory is the factory function to be used to create a GitHub search backend.
func Factory(ctx context.Context, conf *search.BackendConfig) (search.Backend, error) {
	if len(conf.AuthMethod) == 0 {
		// auth method required for GitHub code search - https://developer.github.com/v3/search/#search-code
		if conf.SearchMethod == search.Code {
			return nil, fmt.Errorf("auth method required for code search")
		}
	}

	var entity searchMethodEntity
	if conf.SearchMethod == search.Version {
		if conf.Config["version_entity"] == "" {
			return nil, fmt.Errorf("version_entity required with Version search method")
		}

		entity = conf.Config["version_entity"]
	}

	if conf.SearchMethod == search.Meta {
		if conf.Config["meta_entity"] == "" {
			return nil, fmt.Errorf("meta_entity required with Meta search method")
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

	c := github.NewClient(conf.Client)

	return &Backend{
		auth: auth,
		githubClient: Client{
			SearchService:     c.Search,
			RepositoryService: c.Repositories,
		},
		searchMethod:       conf.SearchMethod,
		searchMethodEntity: entity,
		maxPageSize:        maxPageSize,
	}, nil
}

// Backend is a GitHub search backend.
type Backend struct {
	auth               transport.AuthMethod
	githubClient       Client
	searchMethod       search.Method
	searchMethodEntity searchMethodEntity
	maxPageSize        int
}

// Search is the search function for searching GitHub for projects, code snippets,
// labels, topics, etc. and transparently paginating results.
func (b *Backend) Search(ctx context.Context, query string, numDesiredResults int) ([]project.Backend, error) {
	res := make([]project.Backend, 0, numDesiredResults)

	var page int
	opts := github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: pageSize(numDesiredResults, b.maxPageSize),
		},
	}

	for {
		var (
			searchRes []project.Backend
			resp      *github.Response
			err       error
		)

		switch b.searchMethod {
		case search.Project:
			searchRes, resp, err = searchRepositories(ctx, b.githubClient, query, numDesiredResults, &opts)
		case search.Code:
			searchRes, resp, err = searchCode(ctx, b.githubClient, query, numDesiredResults, &opts)
		case search.Meta:
			searchRes, resp, err = searchMeta(ctx, b.searchMethodEntity, &opts)
		default:
			return nil, fmt.Errorf("unsupported search method")
		}

		if err != nil {
			return nil, err
		}

		for _, r := range searchRes {
			res = append(res, r)
			if len(res) >= numDesiredResults {
				return res, nil
			}
		}

		if resp == nil || resp.NextPage == 0 {
			return res, ErrFewerResultsThanDesired
		}

		opts.ListOptions.Page = page + 1
	}
}
