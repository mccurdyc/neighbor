package binary

import (
	"context"
	"fmt"
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
		"empty_cmd": {
			input: input{
				&run.BackendConfig{
					Cmd: "",
				},
			},
			want: want{
				backend: nil,
				err:     fmt.Errorf("command cannot be nil"),
			},
		},

		"executable_not_found": {
			input: input{
				&run.BackendConfig{
					Cmd: "12345678900987654321",
				},
			},
			want: want{
				backend: nil,
				err:     fmt.Errorf(`failed to find command: 'exec: "12345678900987654321": executable file not found in $PATH'`),
			},
		},

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

		"cmd_with_single_arg": {
			input: input{
				&run.BackendConfig{
					Cmd: "ls -al",
				},
			},
			want: want{
				backend: &Backend{
					cmd:  "ls -al",
					name: "ls",
					args: []string{"-al"},
				},
				err: nil,
			},
		},

		"cmd_with_multiple_args": {
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

func Test_Run(t *testing.T) {
	type input struct {
		backend *Backend
		dir     string
	}

	type want struct {
		err error
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"empty_dir": {
			input: input{
				backend: &Backend{name: "ls", args: []string{"-al"}},
				dir:     "",
			},
			want: want{
				err: fmt.Errorf("working directory must be specified"),
			},
		},

		"cmd_no_args": {
			input: input{
				backend: &Backend{name: "ls"},
				dir:     "./testdata/",
			},
			want: want{
				err: nil,
			},
		},

		"cmd_with_one_arg": {
			input: input{
				backend: &Backend{name: "ls", args: []string{"-al"}},
				dir:     "./testdata/",
			},
			want: want{
				err: nil,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := tt.input.backend.Run(context.TODO(), tt.input.dir)

			// https://github.com/google/go-cmp/issues/24
			errorCmp := func(x, y error) bool {
				if x == nil || y == nil {
					return x == nil && y == nil
				}
				return x.Error() == y.Error()
			}

			if ok := errorCmp(gotErr, tt.want.err); !ok {
				t.Errorf("Run() \n\tgotErr: '%+v'\n\twantErr: '%+v'", gotErr, tt.want.err)
			}
		})
	}
}
