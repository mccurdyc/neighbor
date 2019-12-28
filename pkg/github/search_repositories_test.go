package github

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/github"
)

func Test_processResults(t *testing.T) {
	rs := &RepositorySearcher{}

	type input struct {
		res  *github.RepositoriesSearchResult
		resp *github.Response
	}

	type want struct {
		Results
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"one_page_three_results": {
			input: input{
				res: &github.RepositoriesSearchResult{
					Total:             ptrToInt(3),
					IncompleteResults: ptrToBool(false),
					Repositories: []github.Repository{
						github.Repository{ID: ptrToInt64(1)},
						github.Repository{ID: ptrToInt64(2)},
						github.Repository{ID: ptrToInt64(3)},
					},
				},
				resp: &github.Response{
					NextPage:  0,
					PrevPage:  0,
					FirstPage: 0,
					LastPage:  0,
				},
			},
			want: want{
				Results: Results{
					Repositories: []*github.Repository{
						&github.Repository{ID: ptrToInt64(1)},
						&github.Repository{ID: ptrToInt64(2)},
						&github.Repository{ID: ptrToInt64(3)},
					},
					response: &github.Response{
						NextPage:  0,
						PrevPage:  0,
						FirstPage: 0,
						LastPage:  0,
					},
				},
			},
		},

		"multiple_pages_three_results": {
			input: input{
				res: &github.RepositoriesSearchResult{
					Total:             ptrToInt(6),
					IncompleteResults: ptrToBool(true),
					Repositories: []github.Repository{
						github.Repository{ID: ptrToInt64(1)},
						github.Repository{ID: ptrToInt64(2)},
						github.Repository{ID: ptrToInt64(3)},
					},
				},
				resp: &github.Response{
					NextPage:  1,
					PrevPage:  0,
					FirstPage: 0,
					LastPage:  1,
				},
			},
			want: want{
				Results: Results{
					Repositories: []*github.Repository{
						&github.Repository{ID: ptrToInt64(1)},
						&github.Repository{ID: ptrToInt64(2)},
						&github.Repository{ID: ptrToInt64(3)},
					},
					response: &github.Response{
						NextPage:  1,
						PrevPage:  0,
						FirstPage: 0,
						LastPage:  1,
					},
				},
			},
		},
	}

	for name, tt := range tests {

		t.Run(name, func(t *testing.T) {
			got := rs.processResults(tt.input.res, tt.input.resp)

			if diff := cmp.Diff(tt.want.Results.Repositories, got.Repositories); diff != "" {
				t.Errorf("processResults() Repositories: mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tt.want.Results.response, got.response, cmp.AllowUnexported()); diff != "" {
				t.Errorf("processResults() github response: mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
