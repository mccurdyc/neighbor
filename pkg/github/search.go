package github

import (
	// stdlib
	"encoding/json"
	"fmt"

	// external
	"github.com/google/go-github/github"

	// internal
	"github.com/mccurdyc/neighbor/pkg/neighbor"
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
func (s *SearchService) Search(ctx *neighbor.Ctx, t string, q []byte, opts *github.SearchOptions) (interface{}, *github.Response) {
	var query interface{}

	switch t {
	case "repository":
		query = &RepositoryQuery{}
		err := json.Unmarshal(q, query)
		if err != nil {
			ctx.Logger.Error("error unmarshalling into RepositoryQuery")
			return nil, nil
		}
		break
	case "code":
		query = &CodeQuery{}
		err := json.Unmarshal(q, query)
		if err != nil {
			ctx.Logger.Error("error unmarshalling into CodeQuery")
			return nil, nil
		}
		break
	default:
		ctx.Logger.Info("query type not accepted")
		return nil, nil
	}

	switch d := query.(type) {
	case *RepositoryQuery:
		qStr := buildQuery(d)
		res, resp, err := s.Client.Search.Repositories(ctx.Context, qStr, opts)
		if err != nil {
			ctx.Logger.Error("error searching for repositories")
		}
		return res, resp
	case *CodeQuery:
		qStr := buildQuery(d)
		res, resp, err := s.Client.Search.Code(ctx.Context, qStr, opts)
		if err != nil {
			ctx.Logger.Error("error searching for code")
		}
		return res, resp
	default:
		ctx.Logger.Infof("query type not found %q", d)
		return nil, nil
	}
}

// buildQuery builds the appropriate search query based on the type of query.
// TODO(B): continue adding parameters to the query
func buildQuery(q interface{}) string {
	switch d := q.(type) {
	case *RepositoryQuery:
		// FIXME: this needs dynamically built
		// return fmt.Sprintf("%s user:%s language:%s stars:%d", d.Other, d.User, d.Language, d.Stars)
		return fmt.Sprintf("%s user:%s", d.Other, d.User)
	case *CodeQuery:
		return fmt.Sprintf("%s file:%s", d.Other, d.File)
	default:
		return ""
	}
}
