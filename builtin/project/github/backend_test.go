package github

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mccurdyc/neighbor/sdk/project"
	"github.com/mccurdyc/neighbor/sdk/retrieval"
)

func Test_Factory(t *testing.T) {
	type input struct {
		conf *project.BackendConfig
	}

	type want struct {
		be  *Backend
		err error
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"config_with_no_name": {
			input: input{
				conf: &project.BackendConfig{
					Name: "",
				},
			},
			want: want{
				be:  nil,
				err: fmt.Errorf("name cannot be empty"),
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

func compareBackend(t *testing.T, want *Backend, got project.Backend) {
	t.Helper()

	if got == nil {
		if want != nil {
			t.Errorf("Factory() mismatched nil")
		}
		return
	}

	gotProjectBackend, ok := got.(*Backend)
	if !ok {
		t.Errorf("Factory() failed to type convert to project.Backend")
	}

	if diff := cmp.Diff(want.name, gotProjectBackend.name, cmp.AllowUnexported()); diff != "" {
		t.Errorf("Factory() mismatched name (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(want.retrievalFunc, gotProjectBackend.retrievalFunc, cmp.AllowUnexported()); diff != "" {
		t.Errorf("Factory() mismatched retrieval function (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(want.version, gotProjectBackend.version, cmp.AllowUnexported()); diff != "" {
		t.Errorf("Factory() mismatched version (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(want.sourceLocation, gotProjectBackend.sourceLocation, cmp.AllowUnexported()); diff != "" {
		t.Errorf("Factory() mismatched sourceLocation(-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(want.localLocation, gotProjectBackend.localLocation, cmp.AllowUnexported()); diff != "" {
		t.Errorf("Factory() mismatched localLocation(-want +got):\n%s", diff)
	}
}

func Test_Name(t *testing.T) {
	type input struct {
		backend *Backend
	}

	type want struct {
		value string
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"value_set": {
			input: input{
				backend: &Backend{
					name: "here",
				},
			},
			want: want{
				value: "here",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.input.backend.Name()

			if got != tt.want.value {
				t.Errorf("Name(%+v): \n\tgot: '%+v'\n\twant: '%+v'", tt.input, got, tt.want.value)
			}
		})
	}
}

func Test_Version(t *testing.T) {
	type input struct {
		backend *Backend
	}

	type want struct {
		value string
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"value_set": {
			input: input{
				backend: &Backend{
					version: "here",
				},
			},
			want: want{
				value: "here",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.input.backend.Version()

			if got != tt.want.value {
				t.Errorf("Version(%+v): \n\tgot: '%+v'\n\twant: '%+v'", tt.input, got, tt.want.value)
			}
		})
	}
}

type mockRetrievalBackend struct{}

func (m *mockRetrievalBackend) Retrieve(_ context.Context, _ string, _ string) error { return nil }

func Test_RetrievalFunc(t *testing.T) {
	type input struct {
		backend *Backend
	}

	type want struct {
		value retrieval.Backend
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"value_set": {
			input: input{
				backend: &Backend{
					retrievalFunc: &mockRetrievalBackend{},
				},
			},
			want: want{
				value: &mockRetrievalBackend{},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.input.backend.Version()

			if reflect.DeepEqual(got, tt.want.value) {
				t.Errorf("RetrievalFunc(%+v): \n\tgot: '%+v'\n\twant: '%+v'", tt.input, got, tt.want.value)
			}
		})
	}
}

func Test_SourceLocation(t *testing.T) {
	type input struct {
		backend *Backend
	}

	type want struct {
		value string
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"value_set": {
			input: input{
				backend: &Backend{
					sourceLocation: "here",
				},
			},
			want: want{
				value: "here",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.input.backend.SourceLocation()

			if got != tt.want.value {
				t.Errorf("SourceLocation(%+v): \n\tgot: '%+v'\n\twant: '%+v'", tt.input, got, tt.want.value)
			}
		})
	}
}

func Test_LocalLocation(t *testing.T) {
	type input struct {
		backend *Backend
	}

	type want struct {
		value string
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"value_set": {
			input: input{
				backend: &Backend{
					localLocation: "here",
				},
			},
			want: want{
				value: "here",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.input.backend.LocalLocation()

			if got != tt.want.value {
				t.Errorf("LocalLocation(%+v): \n\tgot: '%+v'\n\twant: '%+v'", tt.input, got, tt.want.value)
			}
		})
	}
}

func Test_SetLocalLocation(t *testing.T) {
	type input struct {
		backend *Backend
		l       string
	}

	type want struct {
		value string
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"value_not_set": {
			input: input{
				backend: &Backend{},
				l:       "here",
			},
			want: want{
				value: "here",
			},
		},

		"value_already_set": {
			input: input{
				backend: &Backend{
					localLocation: "here",
				},
				l: "there",
			},
			want: want{
				value: "there",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.input.backend.SetLocalLocation(tt.input.l)
			got := tt.input.backend.LocalLocation()

			if got != tt.want.value {
				t.Errorf("SetLocalLocation(%+v): \n\tgot: '%+v'\n\twant: '%+v'", tt.input, got, tt.want.value)
			}
		})
	}
}
