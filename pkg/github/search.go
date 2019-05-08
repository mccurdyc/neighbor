package github

import (
	// stdlib

	"strings"

	// external
	"github.com/golang/glog"
	"github.com/google/go-github/github"

	// internal
	"github.com/mccurdyc/neighbor/pkg/neighbor"
)

// RepositoryQuery contains all of the supported fields in a GitHub repository query
// References:
//		+ GitHub API Docs: https://developer.github.com/v3/search/#search-repositories
//		+ GitHub Search Repository Docs: https://help.github.com/articles/searching-for-repositories/
// TODO(A): add the additional supported fields with the appropriate types
// for now, we will leave this, but they will be replaced by the GitHub query string
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
// for now, we will leave this, but they will be replaced by the GitHub query string
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
func (s *SearchService) Search(ctx *neighbor.Ctx, t string, q string, opts *github.SearchOptions) (interface{}, *github.Response) {
	glog.V(2).Infof("performing GitHub search with query: %s", q)

	switch strings.ToLower(t) {
	case "repository":
		// TODO: do pagination on resp
		// API Reference: https://developer.github.com/v3/search/
		// Find repositories via various criteria. This method returns up to 100 results per page.
		res, resp, err := s.Client.Search.Repositories(ctx.Context, q, opts)
		if err != nil {
			glog.Errorf("error searching for repositories: %+v", err)
		}
		return res, resp
	case "code":
		res, resp, err := s.Client.Search.Code(ctx.Context, q, opts)
		if err != nil {
			glog.Errorf("error searching for code: %+v", err)
		}
		return res, resp
	default:
		glog.Errorf("query type \"%s\" not accepted", t)
		return nil, nil
	}
}
