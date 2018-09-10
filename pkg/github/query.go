package github

// RepositoryQuery contains all of the supported fields in a GitHub repository query
// References:
//		+ GitHub API Docs: https://developer.github.com/v3/search/#search-repositories
//		+ GitHub Search Repository Docs: https://help.github.com/articles/searching-for-repositories/
// TODO: add the additional supported fields with the appropriate types
type RepositoryQuery struct {
	User     string `json:"user"`
	Language string `json:"language"`
	Stars    int32  `json:"stars"`
}

// CodeQuery contains all of the supported fields in a GitHub code query.
// References:
//		+ GitHub API Docs: https://developer.github.com/v3/search/#search-code
//		+ GitHub Search Repository Docs: https://help.github.com/articles/searching-code/
// TODO: add the additional supported fields with the appropriate types
type CodeQuery struct {
	File string `json:"file"`
}
