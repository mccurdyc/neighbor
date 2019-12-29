package runner

import (
	"errors"
	"os"
	"testing"
)

type runnerEvaluator interface {
	evalDir(*testing.T, string)

	Runner
}

type mockRunner struct {
	err       error
	actualDir string
}

func (mr *mockRunner) Run() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	mr.actualDir = dir
	return mr.err
}

func (mr *mockRunner) evalDir(t *testing.T, expectedDir string) {
	if mr.actualDir != expectedDir {
		t.Errorf("mockRunner eval(tt.input) = '%s', want '%s'", mr.actualDir, expectedDir)
	}
}

func Test_RunInDir(t *testing.T) {
	type input struct {
		runnerEvaluator runnerEvaluator
	}

	type want struct {
		err error
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"runner_returns_an_error": {
			input: input{
				runnerEvaluator: &mockRunner{err: errors.New("mockRunner error")},
			},
			want: want{
				errors.New("mockRunner error"),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			dir := os.TempDir()

			got := RunInDir(dir, tt.input.runnerEvaluator)

			tt.input.runnerEvaluator.evalDir(t, dir)

			// https://github.com/google/go-cmp/issues/24
			errorCmp := func(x, y error) bool {
				if x == nil || y == nil {
					return x == nil && y == nil
				}
				return x.Error() == y.Error()
			}

			if ok := errorCmp(got, tt.want.err); !ok {
				t.Errorf("RunInDir(tt.input) = '%+v', wantErr '%+v'", got, tt.want.err)
			}
		})
	}
}
