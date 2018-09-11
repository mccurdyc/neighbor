package github

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-github/github"
)

// RepositoryQuery contains all of the supported fields in a GitHub repository query
// References:
//		+ GitHub API Docs: https://developer.github.com/v3/search/#search-repositories
//		+ GitHub Search Repository Docs: https://help.github.com/articles/searching-for-repositories/
// TODO(A): add the additional supported fields with the appropriate types
type RepositoryQuery struct {
	Other string `json:"other"`

	User     string `json:"user"`
	Language string `json:"language"`
	Stars    int32  `json:"stars"`
}

// CodeQuery contains all of the supported fields in a GitHub code query.
// References:
//		+ GitHub API Docs: https://developer.github.com/v3/search/#search-code
//		+ GitHub Search Repository Docs: https://help.github.com/articles/searching-code/
// TODO(A): add the additional supported fields with the appropriate types
type CodeQuery struct {
	Other string `json:"other"`

	File string `json:"file"`
}

// SearchService contains the GitHub client and query necessary to query GitHub for arbitrary data.
type SearchService struct {
	Client *github.Client
}

// NewSearchService is a constructor that returns a pointer to a new SearchService object
// with the GitHub client set.
func NewSearchService(c *github.Client) *SearchService {
	return &SearchService{
		Client: c,
	}
}

// Search is a wrapper for the GitHub library search functionality, but where we can
// build the search queries.
// TODO(D): continue adding other search options
func (s *SearchService) Search(ctx context.Context, t string, q []byte, opts *github.SearchOptions) (interface{}, *github.Response, error) {
	var query interface{}

	switch t {
	case "repository":
		query = &RepositoryQuery{}
		err := json.Unmarshal(q, query)
		if err != nil {
			return nil, nil, err
		}
		break
	case "code":
		query = &CodeQuery{}
		err := json.Unmarshal(q, query)
		if err != nil {
			return nil, nil, err
		}
		break
	default:
		return nil, nil, fmt.Errorf("search type not accepted %s", t)
	}

	switch d := query.(type) {
	case *RepositoryQuery:
		qStr := buildQuery(d)
		return s.Client.Search.Repositories(ctx, qStr, opts)
	case *CodeQuery:
		qStr := buildQuery(d)
		return s.Client.Search.Code(ctx, qStr, opts)
	default:
		return nil, nil, fmt.Errorf("query type not found %q", d)
	}
}

// buildQuery builds the appropriate search query based on the type of query.
// TODO(B): continue adding parameters to the query
func buildQuery(q interface{}) string {
	switch d := q.(type) {
	case *RepositoryQuery:
		return fmt.Sprintf("%s user:%s language:%s stars:%d", d.Other, d.User, d.Language, d.Stars)
	case *CodeQuery:
		return fmt.Sprintf("%s file:%s", d.Other, d.File)
	default:
		return ""
	}
}
