package binary

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mccurdyc/neighbor/sdk/run"
)

func Test_Factory(t *testing.T) {
	type input struct {
		conf *run.BackendConfig
	}

	type want struct {
		backend *Backend
		err     error
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"cmd_in_path": {
			input: input{
				&run.BackendConfig{
					Cmd: "go",
				},
			},
			want: want{
				backend: &Backend{
					cmd:  "go",
					name: "go",
					args: nil,
				},
				err: nil,
			},
		},

		"cmd_with_args": {
			input: input{
				&run.BackendConfig{
					Cmd: "go test $(go list ./...) | grep \"FAIL\"",
				},
			},
			want: want{
				backend: &Backend{
					cmd:  "go test $(go list ./...) | grep \"FAIL\"",
					name: "go",
					args: []string{"test", "$(go", "list", "./...)", "|", "grep", "\"FAIL\""},
				},
				err: nil,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, gotErr := Factory(context.TODO(), tt.input.conf)

			compareBackend(t, tt.want.backend, got)

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

func compareBackend(t *testing.T, want *Backend, got run.Backend) {
	t.Helper()

	if got == nil {
		if want != nil {
			t.Errorf("Factory() mismatched nil")
		}
		return
	}

	gotBackend, ok := got.(*Backend)
	if !ok {
		t.Errorf("Factory() failed to type convert to binary.Backend")
	}

	if diff := cmp.Diff(gotBackend.name, want.name, cmp.AllowUnexported()); diff != "" {
		t.Errorf("Factory() mismatched name (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(gotBackend.cmd, want.cmd, cmp.AllowUnexported()); diff != "" {
		t.Errorf("Factory() mismatched cmd (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(gotBackend.args, want.args, cmp.AllowUnexported()); diff != "" {
		t.Errorf("Factory() mismatched args (-want +got):\n%s", diff)
	}
}
