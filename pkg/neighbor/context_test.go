package neighbor

import (
	"testing"
	// external
	// internal
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name        string
		ctx         *Ctx
		expectedErr string
	}{
		{
			"empty-context",
			&Ctx{},
			"search type, query, external command cannot be empty",
		},
		{
			"empty-access-token",
			&Ctx{
				GitHub: GitHubDetails{
					AccessToken: "",
					SearchType:  "repository",
					Query:       "test query",
				},
				ExternalCmd: []string{"ls", "-al"},
			},
			"",
		},
		{
			"empty-search-type",
			&Ctx{
				GitHub: GitHubDetails{
					AccessToken: "abc123",
					SearchType:  "",
					Query:       "test query",
				},
				ExternalCmd: []string{"ls", "-al"},
			},
			"search type cannot be empty",
		},
		{
			"empty-github-query",
			&Ctx{
				GitHub: GitHubDetails{
					AccessToken: "abc123",
					SearchType:  "repository",
					Query:       "",
				},
				ExternalCmd: []string{"ls", "-al"},
			},
			"query cannot be empty",
		},
		{
			"empty-external-command",
			&Ctx{
				GitHub: GitHubDetails{
					AccessToken: "abc123",
					SearchType:  "repository",
					Query:       "test query",
				},
				ExternalCmd: []string{},
			},
			"external command cannot be empty",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.ctx.Validate()

			if actual == nil {
				if c.expectedErr != "" {
					t.Errorf("\n\tACTUAL: %+v\n\tEXPECTED: %+v\n", actual, c.expectedErr)
				}
				return
			}

			if actual.Error() != c.expectedErr {
				t.Errorf("\n\tACTUAL: %+v\n\tEXPECTED: %+v\n", actual.Error(), c.expectedErr)
			}
		})
	}
}
