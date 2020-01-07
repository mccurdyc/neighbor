package github

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/github"
	githubProject "github.com/mccurdyc/neighbor/builtin/project/github"
	"github.com/mccurdyc/neighbor/sdk/project"
	"github.com/mccurdyc/neighbor/sdk/search"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func Test_Factory(t *testing.T) {
	type input struct {
		conf *search.BackendConfig
	}

	type want struct {
		backend *Backend
		err     error
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"missing_auth_code_search_method": {
			input: input{
				conf: &search.BackendConfig{
					AuthMethod:   "",
					SearchMethod: search.Code,
				},
			},
			want: want{
				err: fmt.Errorf("auth method required for code search"),
			},
		},

		"missing_version_entity_version_search": {
			input: input{
				conf: &search.BackendConfig{
					SearchMethod: search.Version,
				},
			},
			want: want{
				err: fmt.Errorf("version_entity required with Version search method"),
			},
		},

		"version_search": {
			input: input{
				conf: &search.BackendConfig{
					SearchMethod: search.Version,
					Config:       map[string]string{"version_entity": "commit"},
				},
			},
			want: want{
				backend: &Backend{
					searchMethod:       search.Version,
					searchMethodEntity: commit,
				},
				err: nil,
			},
		},

		"missing_meta_entity_meta_search": {
			input: input{
				conf: &search.BackendConfig{
					SearchMethod: search.Meta,
				},
			},
			want: want{
				err: fmt.Errorf("meta_entity required with Meta search method"),
			},
		},

		"meta_search": {
			input: input{
				conf: &search.BackendConfig{
					SearchMethod: search.Meta,
					Config:       map[string]string{"meta_entity": "topic"},
				},
			},
			want: want{
				backend: &Backend{
					searchMethod:       search.Meta,
					searchMethodEntity: topic,
				},
				err: nil,
			},
		},

		"missing_username_basic_auth": {
			input: input{
				conf: &search.BackendConfig{
					SearchMethod: search.Code,
					AuthMethod:   "basic",
					Config:       map[string]string{"password": "password123"},
				},
			},
			want: want{
				err: fmt.Errorf("username required for basic auth"),
			},
		},

		"missing_password_basic_auth": {
			input: input{
				conf: &search.BackendConfig{
					SearchMethod: search.Code,
					AuthMethod:   "basic",
					Config:       map[string]string{"username": "username123"},
				},
			},
			want: want{
				err: fmt.Errorf("password required for basic auth"),
			},
		},

		"missing_token_token_auth": {
			input: input{
				conf: &search.BackendConfig{
					SearchMethod: search.Code,
					AuthMethod:   "token",
				},
			},
			want: want{
				err: fmt.Errorf("token required for token auth"),
			},
		},

		"basic_auth": {
			input: input{
				conf: &search.BackendConfig{
					SearchMethod: search.Code,
					AuthMethod:   "basic",
					Config:       map[string]string{"username": "username123", "password": "password123"},
				},
			},
			want: want{
				backend: &Backend{
					searchMethod: search.Code,
					auth: &http.BasicAuth{
						Username: "username123",
						Password: "password123",
					},
				},
				err: nil,
			},
		},

		"token_auth": {
			input: input{
				conf: &search.BackendConfig{
					SearchMethod: search.Code,
					AuthMethod:   "token",
					Config:       map[string]string{"token": "token123"},
				},
			},
			want: want{
				backend: &Backend{
					searchMethod: search.Code,
					auth: &http.BasicAuth{
						Username: "null",
						Password: "token123",
					},
				},
				err: nil,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, gotErr := Factory(context.TODO(), tt.input.conf)

			compareBackend(t, "Factory", tt.want.backend, got)

			// https://github.com/google/go-cmp/issues/24
			errorCmp := func(x, y error) bool {
				if x == nil || y == nil {
					return x == nil && y == nil
				}
				return x.Error() == y.Error()
			}

			if ok := errorCmp(gotErr, tt.want.err); !ok {
				t.Errorf("Factory() \n\tgotErr: '%+v'\n\twantErr: '%+v'", gotErr, tt.want.err)
			}
		})
	}
}

type mockClient struct {
	repositories *github.RepositoriesSearchResult
	commits      []*github.RepositoryCommit
	response     *github.Response
	err          error
}

// Repositories returns the repositories for a given search query.
func (m *mockClient) Repositories(ctx context.Context, query string, opts *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error) {
	return m.repositories, m.response, m.err
}

// ListCommits lists the commits for a specific repository.
func (m *mockClient) ListCommits(ctx context.Context, owner string, repo string, opts *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error) {
	return m.commits, m.response, m.err
}

func newMockClient(maxPageSize int, numCommits int, nextPage bool, err error) Client {
	repos := make([]github.Repository, 0, maxPageSize)

	for i := 0; i < maxPageSize; i++ {
		name := strconv.Itoa(i)
		fullname := fmt.Sprintf("repo/%s", name)
		cloneURL := fmt.Sprintf("cloneurl%d.git", i)
		ownerName := fmt.Sprintf("owner%d", i)

		repos = append(repos, github.Repository{
			Name:     &name,
			FullName: &fullname,
			CloneURL: &cloneURL,
			Owner: &github.User{
				Name: &ownerName,
			},
		})
	}

	commits := make([]*github.RepositoryCommit, 0, numCommits)
	for i := 0; i < numCommits; i++ {
		sha := fmt.Sprintf("sha%d", i)

		commits = append(commits,
			&github.RepositoryCommit{
				SHA: &sha,
			})
	}

	var resp github.Response
	if nextPage {
		resp = github.Response{
			NextPage: 1,
		}
	}

	return Client{
		RepositoryService: &mockClient{
			commits:  commits,
			response: &github.Response{},
			err:      err,
		},
		SearchService: &mockClient{
			repositories: &github.RepositoriesSearchResult{
				Repositories: repos,
			},
			response: &resp,
			err:      err,
		},
	}
}

func generateWantProjects(t *testing.T, name string, numDesiredResults int, maxPageSize int, nextPage bool) []project.Backend {
	t.Helper()

	wantProjects := make([]project.Backend, 0, numDesiredResults)
	for len(wantProjects) < numDesiredResults {
		for i := 0; i < maxPageSize; i++ {
			if len(wantProjects) >= numDesiredResults {
				break
			}

			project, err := githubProject.Factory(context.TODO(), &project.BackendConfig{
				Name:           fmt.Sprintf("repo/%d", i),
				Version:        "sha0", // we always want the latest commit
				SourceLocation: fmt.Sprintf("cloneurl%d.git", i),
			})
			if err != nil {
				t.Errorf("%s(): failed to create project: %+v", name, err)
			}

			wantProjects = append(wantProjects, project)
		}

		if nextPage == false {
			break
		}
	}

	return wantProjects
}

func Test_Search(t *testing.T) {
	type input struct {
		backend           *Backend
		numDesiredResults int
		numCommits        int
		maxPageSize       int
		nextPage          bool
		clientErr         error
	}

	type want struct {
		projectsFn func(*testing.T, string, int, int, bool) []project.Backend
		err        error
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"numDesiredResults_equals_maxPageSize": {
			input: input{
				backend: &Backend{
					searchMethod: search.Project,
				},
				maxPageSize:       3,
				numDesiredResults: 3,
				numCommits:        2,
				clientErr:         nil,
				nextPage:          false,
			},
			want: want{
				projectsFn: generateWantProjects,
				err:        nil,
			},
		},

		"numDesiredResults_greater_than_maxPageSize": {
			input: input{
				backend: &Backend{
					searchMethod: search.Project,
				},
				maxPageSize:       3,
				numDesiredResults: 5,
				numCommits:        2,
				clientErr:         nil,
				nextPage:          true,
			},
			want: want{
				projectsFn: generateWantProjects,
				err:        nil,
			},
		},

		"less_than_desired": {
			input: input{
				backend: &Backend{
					searchMethod: search.Project,
				},
				maxPageSize:       3,
				numDesiredResults: 5,
				numCommits:        2,
				clientErr:         nil,
				nextPage:          false,
			},
			want: want{
				projectsFn: generateWantProjects,
				err:        ErrFewerResultsThanDesired,
			},
		},

		"github_client_error": {
			input: input{
				backend: &Backend{
					searchMethod: search.Project,
				},
				maxPageSize:       3,
				numDesiredResults: 5,
				numCommits:        2,
				clientErr:         fmt.Errorf("github client error"),
				nextPage:          false,
			},
			want: want{
				projectsFn: func(_ *testing.T, _ string, _ int, _ int, _ bool) []project.Backend { return nil },
				err:        fmt.Errorf("github client error"),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.input.backend.githubClient = newMockClient(tt.input.maxPageSize, tt.input.numCommits, tt.input.nextPage, tt.input.clientErr)

			if tt.want.projectsFn == nil {
				t.Fatal("Search(): a want project generation function is required")
			}
			wantProjects := tt.want.projectsFn(t, "Search", tt.input.numDesiredResults, tt.input.maxPageSize, tt.input.nextPage)

			got, gotErr := tt.input.backend.Search(context.TODO(), "query", tt.input.numDesiredResults)

			// https://github.com/google/go-cmp/issues/24
			errorCmp := func(x, y error) bool {
				if x == nil || y == nil {
					return x == nil && y == nil
				}
				return x.Error() == y.Error()
			}

			if ok := errorCmp(gotErr, tt.want.err); !ok {
				t.Errorf("Search() \n\tgotErr: '%+v'\n\twantErr: '%+v'", gotErr, tt.want.err)
			}

			if len(wantProjects) != len(got) {
				t.Errorf("Search() returned a different amount of results: \n\twant: %+v\n\tgot: %+v", len(wantProjects), len(got))
				t.FailNow()
			}

			for i := 0; i < len(wantProjects); i++ {
				compareProject(t, "Search", wantProjects[i], got[i])
			}
		})
	}
}

func compareProject(t *testing.T, name string, want, got project.Backend) {
	t.Helper()

	if want == nil && got != nil {
		t.Errorf("%s(): mismatched nil projects", name)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("%s(): mismatched\n\twant: %+v\n\tgot: %+v", name, want, got)
	}
}

func compareBackend(t *testing.T, name string, want *Backend, got search.Backend) {
	t.Helper()

	if got == nil {
		if want != nil {
			t.Errorf("Factory() mismatched nil")
		}
		return
	}

	gotBackend, ok := got.(*Backend)
	if !ok {
		t.Errorf("Factory() failed to type convert to search.Backend")
	}

	if diff := cmp.Diff(want.auth, gotBackend.auth, cmp.AllowUnexported()); diff != "" {
		t.Errorf("Factory() mismatched auth (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(want.searchMethod, gotBackend.searchMethod, cmp.AllowUnexported()); diff != "" {
		t.Errorf("Factory() mismatched search method (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(want.searchMethodEntity, gotBackend.searchMethodEntity, cmp.AllowUnexported()); diff != "" {
		t.Errorf("Factory() mismatched search method entity (-want +got):\n%s", diff)
	}
}
