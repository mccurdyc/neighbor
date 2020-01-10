package github

import "testing"

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
