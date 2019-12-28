package github

import (
	"context"
	"errors"
	"net/url"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/github"
)

func ptrToInt(i int) *int {
	return &i
}

func ptrToInt64(i int64) *int64 {
	return &i
}

func ptrToBool(b bool) *bool {
	return &b
}

func Test_SearchOptions(t *testing.T) {
	want := searchOptions{
		numDesiredResults:   maxPageSize,
		maxPageSize:         maxPageSize,
		gitHubSearchOptions: github.SearchOptions{},
	}

	got := SearchOptions()
	eq := reflect.DeepEqual(want, got)

	if !eq {
		t.Errorf("SearchOptions() mismatch \n\tgot: %+v\n\twant: %+v", got, want)
	}
}

func Test_WitNumberOfResults(t *testing.T) {
	type input struct {
		so searchOptions
		n  int
	}

	type want struct {
		so searchOptions
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"empty_search_options": {
			input: input{
				so: searchOptions{numDesiredResults: 0, maxPageSize: 0, gitHubSearchOptions: github.SearchOptions{}},
				n:  1,
			},
			want: want{
				so: searchOptions{numDesiredResults: 1, maxPageSize: 0, gitHubSearchOptions: github.SearchOptions{}},
			},
		},

		"non_empty_search_options": {
			input: input{
				so: searchOptions{numDesiredResults: 1, maxPageSize: 1, gitHubSearchOptions: github.SearchOptions{
					Sort:      "sort",
					Order:     "order",
					TextMatch: true,
					ListOptions: github.ListOptions{
						Page:    1,
						PerPage: 1,
					},
				}},
				n: 5,
			},
			want: want{
				so: searchOptions{numDesiredResults: 5, maxPageSize: 1, gitHubSearchOptions: github.SearchOptions{
					Sort:      "sort",
					Order:     "order",
					TextMatch: true,
					ListOptions: github.ListOptions{
						Page:    1,
						PerPage: 1,
					},
				}},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.input.so.WithNumberOfResults(tt.input.n)

			if diff := cmp.Diff(tt.want.so.maxPageSize, got.maxPageSize, cmp.AllowUnexported()); diff != "" {
				t.Errorf("WithNumberOfResults() maxPageSize: mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tt.want.so.numDesiredResults, got.numDesiredResults, cmp.AllowUnexported()); diff != "" {
				t.Errorf("WithNumberOfResults() numDesiredResults: mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tt.want.so.gitHubSearchOptions, got.gitHubSearchOptions, cmp.AllowUnexported()); diff != "" {
				t.Errorf("WithNumberOfResults() gitHubSearchOptions: mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_WithGitHubOptions(t *testing.T) {
	type input struct {
		so                  searchOptions
		gitHubSearchOptions github.SearchOptions
	}

	type want struct {
		so searchOptions
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"empty_search_options": {
			input: input{
				so: searchOptions{numDesiredResults: 0, maxPageSize: 0, gitHubSearchOptions: github.SearchOptions{}},
				gitHubSearchOptions: github.SearchOptions{
					Sort:      "sort",
					Order:     "order",
					TextMatch: true,
					ListOptions: github.ListOptions{
						Page:    2,
						PerPage: 2,
					},
				},
			},
			want: want{
				so: searchOptions{numDesiredResults: 0, maxPageSize: 0, gitHubSearchOptions: github.SearchOptions{
					Sort:      "sort",
					Order:     "order",
					TextMatch: true,
					ListOptions: github.ListOptions{
						Page:    2,
						PerPage: 2,
					},
				},
				},
			},
		},

		"non_empty_search_options": {
			input: input{
				so: searchOptions{numDesiredResults: 1, maxPageSize: 1, gitHubSearchOptions: github.SearchOptions{
					Sort:      "old",
					Order:     "old",
					TextMatch: true,
					ListOptions: github.ListOptions{
						Page:    1,
						PerPage: 1,
					},
				},
				},
				gitHubSearchOptions: github.SearchOptions{
					Sort:      "new",
					Order:     "new",
					TextMatch: true,
					ListOptions: github.ListOptions{
						Page:    2,
						PerPage: 2,
					},
				},
			},
			want: want{
				so: searchOptions{numDesiredResults: 1, maxPageSize: 1, gitHubSearchOptions: github.SearchOptions{
					Sort:      "new",
					Order:     "new",
					TextMatch: true,
					ListOptions: github.ListOptions{
						Page:    2,
						PerPage: 2,
					},
				},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.input.so.WithGitHubOptions(tt.input.gitHubSearchOptions)

			if diff := cmp.Diff(tt.want.so.maxPageSize, got.maxPageSize, cmp.AllowUnexported()); diff != "" {
				t.Errorf("WithGitHubOptions() maxPageSize: mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tt.want.so.numDesiredResults, got.numDesiredResults, cmp.AllowUnexported()); diff != "" {
				t.Errorf("WithGitHubOptions() numDesiredResults: mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tt.want.so.gitHubSearchOptions, got.gitHubSearchOptions, cmp.AllowUnexported()); diff != "" {
				t.Errorf("WithGitHubOptions() gitHubSearchOptions: mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_NewSearcher(t *testing.T) {
	type input struct {
		c *github.Client
		t SearchType
	}

	type want struct {
		s   Searcher
		err error
	}

	ghClient := github.Client{BaseURL: &url.URL{Host: "localhost"}}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"unsupported_search_type": {
			input: input{
				c: &ghClient,
				t: SearchType("unsupported"),
			},
			want: want{
				s:   nil,
				err: errors.New("unsupported search type"),
			},
		},

		"code_search_type": {
			input: input{
				c: &ghClient,
				t: SearchType("code"),
			},
			want: want{
				s:   &CodeSearcher{client: &ghClient},
				err: nil,
			},
		},

		"repository_search_type": {
			input: input{
				c: &ghClient,
				t: SearchType("repository"),
			},
			want: want{
				s:   &RepositorySearcher{client: &ghClient},
				err: nil,
			},
		},

		"repository_search_type_improper_casing_returns_error": {
			input: input{
				c: &ghClient,
				t: SearchType("RepositorY"),
			},
			want: want{
				s:   nil,
				err: errors.New("unsupported search type"),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, gotErr := NewSearcher(tt.input.c, tt.input.t)

			if !reflect.DeepEqual(got, tt.want.s) {
				t.Errorf("NewSearcher(%+v): want '%s', got '%s'", tt.input, tt.want.s, got)
			}

			// https://github.com/google/go-cmp/issues/24
			errorCmp := func(x, y error) bool {
				if x == nil || y == nil {
					return x == nil && y == nil
				}
				return x.Error() == y.Error()
			}

			if ok := errorCmp(gotErr, tt.want.err); !ok {
				t.Errorf("NewSearcher(tt.input) = %v, want %v", gotErr, tt.want.err)
			}
		})
	}
}

type mockSearcher struct {
	res Results
	err error
}

func (m *mockSearcher) search(_ context.Context, _ string, _ *github.SearchOptions) (Results, error) {
	return m.res, m.err
}

func (m *mockSearcher) processResults(_ interface{}, _ *github.Response) Results {
	return m.res
}

func Test_Search(t *testing.T) {
	type input struct {
		s     Searcher
		query string
		opts  searchOptions
	}

	type want struct {
		res []*github.Repository
		err error
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"no_results_returned": {
			input: input{
				s: &mockSearcher{
					res: Results{},
					err: nil,
				},
				query: "",
				opts:  searchOptions{numDesiredResults: 3, maxPageSize: 5},
			},
			want: want{
				res: []*github.Repository{},
				err: ErrFewerResultsThanDesired,
			},
		},

		"single_results_page": {
			input: input{
				s: &mockSearcher{
					res: Results{
						Repositories: []*github.Repository{
							&github.Repository{ID: ptrToInt(1)},
							&github.Repository{ID: ptrToInt(2)},
							&github.Repository{ID: ptrToInt(3)},
							&github.Repository{ID: ptrToInt(4)},
							&github.Repository{ID: ptrToInt(5)},
						},
					},
					err: nil,
				},
				query: "",
				opts:  searchOptions{numDesiredResults: 3, maxPageSize: 5},
			},
			want: want{
				res: []*github.Repository{
					&github.Repository{ID: ptrToInt(1)},
					&github.Repository{ID: ptrToInt(2)},
					&github.Repository{ID: ptrToInt(3)},
				},
				err: nil,
			},
		},

		"two_results_pages": {
			input: input{
				s: &mockSearcher{
					res: Results{
						Repositories: []*github.Repository{
							&github.Repository{ID: ptrToInt(1)},
							&github.Repository{ID: ptrToInt(2)},
							&github.Repository{ID: ptrToInt(3)},
						},
						response: &github.Response{NextPage: 1},
					},
					err: nil,
				},
				query: "",
				opts:  searchOptions{numDesiredResults: 5, maxPageSize: 3},
			},
			want: want{
				res: []*github.Repository{
					&github.Repository{ID: ptrToInt(1)},
					&github.Repository{ID: ptrToInt(2)},
					&github.Repository{ID: ptrToInt(3)},
					&github.Repository{ID: ptrToInt(1)},
					&github.Repository{ID: ptrToInt(2)},
				},
				err: nil,
			},
		},

		"three_results_pages": {
			input: input{
				s: &mockSearcher{
					res: Results{
						Repositories: []*github.Repository{
							&github.Repository{ID: ptrToInt(1)},
							&github.Repository{ID: ptrToInt(2)},
						},
						response: &github.Response{NextPage: 1},
					},
					err: nil,
				},
				query: "",
				opts:  searchOptions{numDesiredResults: 5, maxPageSize: 2},
			},
			want: want{
				res: []*github.Repository{
					&github.Repository{ID: ptrToInt(1)},
					&github.Repository{ID: ptrToInt(2)},
					&github.Repository{ID: ptrToInt(1)},
					&github.Repository{ID: ptrToInt(2)},
					&github.Repository{ID: ptrToInt(1)},
				},
				err: nil,
			},
		},

		"search_returns_error": {
			input: input{
				s: &mockSearcher{
					res: Results{
						Repositories: []*github.Repository{
							&github.Repository{ID: ptrToInt(1)},
							&github.Repository{ID: ptrToInt(2)},
						},
						response: &github.Response{NextPage: 1},
					},
					err: errors.New("search error"),
				},
				query: "",
				opts:  searchOptions{numDesiredResults: 5, maxPageSize: 2},
			},
			want: want{
				res: []*github.Repository{},
				err: errors.New("search error"),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			got, gotErr := Search(context.TODO(), tt.input.s, tt.input.query, tt.input.opts)

			if diff := cmp.Diff(tt.want.res, got); diff != "" {
				t.Errorf("Search() mismatch (-want +got):\n%s", diff)
			}

			// https://github.com/google/go-cmp/issues/24
			errorCmp := func(x, y error) bool {
				if x == nil || y == nil {
					return x == nil && y == nil
				}
				return x.Error() == y.Error()
			}

			if ok := errorCmp(gotErr, tt.want.err); !ok {
				t.Errorf("Search(%+v) = %v, wantErr %v", tt.input, gotErr, tt.want.err)
			}
		})
	}
}

func Test_pageSize(t *testing.T) {
	type input struct {
		desired int
		max     int
	}

	type want struct {
		value int
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"desired_greater_than": {
			input: input{
				desired: 5,
				max:     4,
			},
			want: want{
				value: 4,
			},
		},

		"desired_less_than": {
			input: input{
				desired: 3,
				max:     4,
			},
			want: want{
				value: 3,
			},
		},

		"desired_equal_to": {
			input: input{
				desired: 4,
				max:     4,
			},
			want: want{
				value: 4,
			},
		},
	}

	for name, tt := range tests {

		t.Run(name, func(t *testing.T) {
			got := pageSize(tt.input.desired, tt.input.max)

			if got != tt.want.value {
				t.Errorf("pageSize(%+v): want %+v, got %+v", tt.input, tt.want.value, got)
			}
		})
	}
}
