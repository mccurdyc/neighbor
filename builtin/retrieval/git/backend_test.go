package git

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mccurdyc/neighbor/sdk/retrieval"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func Test_Factory(t *testing.T) {
	type input struct {
		conf *retrieval.BackendConfig
	}

	type want struct {
		be  *Backend
		err error
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"config_with_no_auth": {
			input: input{
				conf: &retrieval.BackendConfig{
					AuthMethod: "",
				},
			},
			want: want{
				be: &Backend{
					auth: nil,
				},
				err: nil,
			},
		},

		"config_with_basic_auth_username_password": {
			input: input{
				conf: &retrieval.BackendConfig{
					AuthMethod: "basic",
					Config:     map[string]string{"username": "abc123", "password": "passabc123"},
				},
			},
			want: want{
				be: &Backend{
					auth: &http.BasicAuth{Username: "abc123", Password: "passabc123"},
				},
				err: nil,
			},
		},

		"config_with_basic_auth_missing_username": {
			input: input{
				conf: &retrieval.BackendConfig{
					AuthMethod: "basic",
					Config:     map[string]string{"password": "passabc123"},
				},
			},
			want: want{
				be:  nil,
				err: fmt.Errorf("username required for basic auth"),
			},
		},

		"config_with_basic_auth_missing_password": {
			input: input{
				conf: &retrieval.BackendConfig{
					AuthMethod: "basic",
					Config:     map[string]string{"username": "abc123"},
				},
			},
			want: want{
				be:  nil,
				err: fmt.Errorf("password required for basic auth"),
			},
		},

		"config_with_token_auth": {
			input: input{
				conf: &retrieval.BackendConfig{
					AuthMethod: "token",
					Config:     map[string]string{"token": "abc123"},
				},
			},
			want: want{
				be: &Backend{
					auth: &http.BasicAuth{Username: "null", Password: "abc123"},
				},
				err: nil,
			},
		},

		"config_with_token_auth_missing_token": {
			input: input{
				conf: &retrieval.BackendConfig{
					AuthMethod: "token",
					Config:     map[string]string{},
				},
			},
			want: want{
				be:  nil,
				err: fmt.Errorf("token required for token auth"),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, gotErr := Factory(context.TODO(), tt.input.conf)

			compareBackend(t, tt.want.be, got)

			// https://github.com/google/go-cmp/issues/24
			errorCmp := func(x, y error) bool {
				if x == nil || y == nil {
					return x == nil && y == nil
				}
				return x.Error() == y.Error()
			}

			if ok := errorCmp(gotErr, tt.want.err); !ok {
				t.Errorf("Factory(): \n\tgotErr: '%v'\n\twantErr: '%v'", gotErr, tt.want.err)
			}
		})
	}
}

func compareBackend(t *testing.T, want *Backend, got retrieval.Backend) {
	if got == nil {
		if want != nil {
			t.Errorf("Factory() mismatched nil")
		}
		return
	}

	gotGitBackend, ok := got.(*Backend)
	if !ok {
		t.Errorf("Factory() failed to type convert to git.Backend")
	}

	if diff := cmp.Diff(want.auth, gotGitBackend.auth, cmp.AllowUnexported()); diff != "" {
		t.Errorf("Factory() mismatched auth (-want +got):\n%s", diff)
	}
}

func Test_Retrieval(t *testing.T) {
}
