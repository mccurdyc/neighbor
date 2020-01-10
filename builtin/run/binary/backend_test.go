package binary

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_cleanWords(t *testing.T) {
	type input struct {
		args []string
	}

	type want struct {
		value []string
	}

	var tests = map[string]struct {
		input input
		want  want
	}{
		"regex no match": {
			input: input{
				args: []string{"hello there", "world"},
			},
			want: want{
				value: []string{"hello there", "world"},
			},
		},

		"regex match double quote": {
			input: input{
				args: []string{"\"hello there\"", "\"world\""},
			},
			want: want{
				value: []string{"hello there", "world"},
			},
		},

		"regex match single quote": {
			input: input{
				args: []string{"'hello there'", "'world'"},
			},
			want: want{
				value: []string{"hello there", "world"},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := cleanWords(tt.input.args)

			if diff := cmp.Diff(tt.want.value, got); diff != "" {
				t.Errorf("cleanWords() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
