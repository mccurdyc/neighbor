package github

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
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
